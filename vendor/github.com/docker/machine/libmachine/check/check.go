package check

import (
	"fmt"
	"net/url"

	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/cert"
	"github.com/docker/machine/libmachine/host"
)

var (
	DefaultConnChecker ConnChecker
)

func init() {
	DefaultConnChecker = &MachineConnChecker{}
}

// ErrCertInvalid for when the cert is computed to be invalid.
type ErrCertInvalid struct {
	wrappedErr error
	hostURL    string
}

func (e ErrCertInvalid) Error() string {
	return fmt.Sprintf(`There was an error validating certificates for host %q: %s
You can attempt to regenerate them using 'lambda-machine-local regenerate-certs [name]'.
Be advised that this will trigger a Docker daemon restart which might stop running containers.
`, e.hostURL, e.wrappedErr)
}

type ConnChecker interface {
	Check(*host.Host, bool) (dockerHost string, authOptions *auth.Options, err error)
}

type MachineConnChecker struct{}

// leave swarm bool here to make golang interface checking happy
func (mcc *MachineConnChecker) Check(h *host.Host, swarm bool) (string, *auth.Options, error) {
	dockerHost, err := h.Driver.GetURL()
	if err != nil {
		return "", &auth.Options{}, err
	}

	dockerURL := dockerHost

	u, err := url.Parse(dockerURL)
	if err != nil {
		return "", &auth.Options{}, fmt.Errorf("Error parsing URL: %s", err)
	}

	authOptions := h.AuthOptions()

	if err := checkCert(u.Host, authOptions); err != nil {
		return "", &auth.Options{}, fmt.Errorf("Error checking and/or regenerating the certs: %s", err)
	}

	return dockerURL, authOptions, nil
}

func checkCert(hostURL string, authOptions *auth.Options) error {
	valid, err := cert.ValidateCertificate(hostURL, authOptions)
	if !valid || err != nil {
		return ErrCertInvalid{
			wrappedErr: err,
			hostURL:    hostURL,
		}
	}

	return nil
}
