package provision

import (
	"github.com/docker/machine/libmachine/drivers"
)

func NewCentosProvisioner(d drivers.Driver) Provisioner {
	return &CentosProvisioner{
		NewRedHatProvisioner("centos", d),
	}
}

type CentosProvisioner struct {
	*RedHatProvisioner
}

func (provisioner *CentosProvisioner) String() string {
	return "centos"
}
