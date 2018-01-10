package provision

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/cert"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/provision/pkgaction"
	"github.com/docker/machine/libmachine/provision/serviceaction"
	"github.com/docker/machine/libmachine/swarm"
)

func init() {
	Register("Lambda Linux", &RegisteredProvisioner{
		New: NewLambdaProvisioner,
	})
}

func NewLambdaProvisioner(d drivers.Driver) Provisioner {
	return &LambdaProvisioner{
		Driver: d,
	}
}

type LambdaProvisioner struct {
	OsReleaseInfo *OsRelease
	Driver        drivers.Driver
	AuthOptions   auth.Options
	EngineOptions engine.Options
}

func (provisioner *LambdaProvisioner) String() string {
	return "lambda"
}

func (provisioner *LambdaProvisioner) Service(name string, action serviceaction.ServiceAction) error {
	_, err := provisioner.SSHCommand(fmt.Sprintf("sudo service %s %s", name, action.String()))
	return err
}

func (provisioner *LambdaProvisioner) CompatibleWithHost() bool {
	return provisioner.OsReleaseInfo.ID == "lambda"
}

func (provisioner *LambdaProvisioner) SetOsReleaseInfo(info *OsRelease) {
	provisioner.OsReleaseInfo = info
}

func (provisioner *LambdaProvisioner) GetOsReleaseInfo() (*OsRelease, error) {
	return provisioner.OsReleaseInfo, nil
}

func (provisioner *LambdaProvisioner) Package(name string, action pkgaction.PackageAction) error {
	return fmt.Errorf("Not supported in Lambda Linux Provisioner")
}

func (provisioner *LambdaProvisioner) Hostname() (string, error) {
	return provisioner.SSHCommand("hostname")
}

func (provisioner *LambdaProvisioner) SetHostname(hostname string) error {
	// We do not set hostname here. The hostname is set using
	// `/etc/dhcp/dhclient.d/sethostname.sh`
	return nil
}

func (provisioner *LambdaProvisioner) GetDockerOptionsDir() string {
	// This is where the TLS certificates will get copied to
	return "/etc/docker/certs"
}

func (provisioner *LambdaProvisioner) GetAuthOptions() auth.Options {
	return provisioner.AuthOptions
}

func (provisioner *LambdaProvisioner) GetSwarmOptions() swarm.Options {
	// Fake swarm
	return swarm.Options{}
}

func (provisioner *LambdaProvisioner) GenerateDockerOptions(dockerPort int) (*DockerOptions, error) {
	// We do not use this to generate docker options. Our docker
	// `/var/lib/lambda-machine-local/etc/sysconfig/docker`
	return &DockerOptions{}, nil
}

func (provisioner *LambdaProvisioner) Provision(swarmOptions swarm.Options, authOptions auth.Options, engineOptions engine.Options) error {
	var (
		err error
	)

	provisioner.AuthOptions = authOptions
	provisioner.EngineOptions = engineOptions

	// stop docker if it is running
	if err := provisioner.Service("docker", serviceaction.Stop); err != nil {
		return err
	}
	// restart lvm2-monitor as a good measure
	if err := provisioner.Service("lvm2-monitor", serviceaction.Restart); err != nil {
		return err
	}

	// We don't use `makeDockerOptionsDir` as that is taken care of by
	// imagebuilder

	provisioner.AuthOptions = setRemoteAuthOptions(provisioner)

	// NOTE: we do not use `ConfigureAuth`. Instead we use `lambdaConfigureAuth`
	if err = lambdaConfigureAuth(provisioner); err != nil {
		return err
	}

	// restart lvm2-monitor as a good measure before starting docker
	if err := provisioner.Service("lvm2-monitor", serviceaction.Restart); err != nil {
		return err
	}
	// start docker
	if err := provisioner.Service("docker", serviceaction.Start); err != nil {
		return err
	}

	return WaitForDocker(provisioner, engine.DefaultPort)
}

func (provisioner *LambdaProvisioner) SSHCommand(args string) (string, error) {
	return drivers.RunSSHCommandFromDriver(provisioner.Driver, args)
}

func (provisioner *LambdaProvisioner) GetDriver() drivers.Driver {
	return provisioner.Driver
}

func lambdaConfigureAuth(p Provisioner) error {
	var (
		err error
	)

	driver := p.GetDriver()
	machineName := driver.GetMachineName()
	authOptions := p.GetAuthOptions()

	org := mcnutils.GetUsername() + "." + machineName
	bits := 2048

	ip, err := driver.GetIP()
	if err != nil {
		return err
	}

	log.Info("Copying certs to the local machine directory...")

	if err := mcnutils.CopyFile(authOptions.CaCertPath, filepath.Join(authOptions.StorePath, "ca.pem")); err != nil {
		return fmt.Errorf("Copying ca.pem to machine dir failed: %s", err)
	}

	if err := mcnutils.CopyFile(authOptions.ClientCertPath, filepath.Join(authOptions.StorePath, "cert.pem")); err != nil {
		return fmt.Errorf("Copying cert.pem to machine dir failed: %s", err)
	}

	if err := mcnutils.CopyFile(authOptions.ClientKeyPath, filepath.Join(authOptions.StorePath, "key.pem")); err != nil {
		return fmt.Errorf("Copying key.pem to machine dir failed: %s", err)
	}

	// The Host IP is always added to the certificate's SANs list
	hosts := append(authOptions.ServerCertSANs, ip, "localhost")
	log.Debugf("generating server cert: %s ca-key=%s private-key=%s org=%s san=%s",
		authOptions.ServerCertPath,
		authOptions.CaCertPath,
		authOptions.CaPrivateKeyPath,
		org,
		hosts,
	)

	// TODO: Switch to passing just authOptions to this func
	// instead of all these individual fields
	err = cert.GenerateCert(&cert.Options{
		Hosts:       hosts,
		CertFile:    authOptions.ServerCertPath,
		KeyFile:     authOptions.ServerKeyPath,
		CAFile:      authOptions.CaCertPath,
		CAKeyFile:   authOptions.CaPrivateKeyPath,
		Org:         org,
		Bits:        bits,
		SwarmMaster: false,
	})

	if err != nil {
		return fmt.Errorf("error generating server cert: %s", err)
	}

	// docker service should be stopped in our case

	if _, err := p.SSHCommand(`if [ ! -z "$(ip link show docker0)" ]; then sudo ip link delete docker0; fi`); err != nil {
		return err
	}

	// upload certs and configure TLS auth
	caCert, err := ioutil.ReadFile(authOptions.CaCertPath)
	if err != nil {
		return err
	}

	serverCert, err := ioutil.ReadFile(authOptions.ServerCertPath)
	if err != nil {
		return err
	}
	serverKey, err := ioutil.ReadFile(authOptions.ServerKeyPath)
	if err != nil {
		return err
	}

	log.Info("Copying certs to the remote machine...")

	// printf will choke if we don't pass a format string because of the
	// dashes, so that's the reason for the '%%s'
	certTransferCmdFmt := "printf '%%s' '%s' | sudo tee %s"

	// These ones are for Jessie and Mike <3 <3 <3
	if _, err := p.SSHCommand(fmt.Sprintf(certTransferCmdFmt, string(caCert), authOptions.CaCertRemotePath)); err != nil {
		return err
	}

	if _, err := p.SSHCommand(fmt.Sprintf(certTransferCmdFmt, string(serverCert), authOptions.ServerCertRemotePath)); err != nil {
		return err
	}

	if _, err := p.SSHCommand(fmt.Sprintf(certTransferCmdFmt, string(serverKey), authOptions.ServerKeyRemotePath)); err != nil {
		return err
	}

	return nil
}
