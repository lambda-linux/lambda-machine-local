package commands

import (
	"github.com/docker/machine/libmachine"
)

func cmdStart(c CommandLine, api libmachine.API) error {
	if err := runAction("start", c, api); err != nil {
		return err
	}

	return nil
}
