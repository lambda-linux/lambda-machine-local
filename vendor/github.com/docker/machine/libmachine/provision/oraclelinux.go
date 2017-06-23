package provision

import (
	"github.com/docker/machine/libmachine/drivers"
)

func NewOracleLinuxProvisioner(d drivers.Driver) Provisioner {
	return &OracleLinuxProvisioner{
		NewRedHatProvisioner("ol", d),
	}
}

type OracleLinuxProvisioner struct {
	*RedHatProvisioner
}

func (provisioner *OracleLinuxProvisioner) String() string {
	return "ol"
}
