package provision

import (
	"fmt"

	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/engine"
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
	return nil
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
	// `/etc/udhcpc/post-bound/set-hostname`
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
	return nil
}

func (provisioner *LambdaProvisioner) SSHCommand(args string) (string, error) {
	return drivers.RunSSHCommandFromDriver(provisioner.Driver, args)
}

func (provisioner *LambdaProvisioner) GetDriver() drivers.Driver {
	return provisioner.Driver
}
