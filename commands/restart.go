package commands

import (
	"github.com/docker/machine/libmachine"
)

func cmdRestart(c CommandLine, api libmachine.API) error {
	if err := runAction("restart", c, api); err != nil {
		return err
	}

	return nil
}
