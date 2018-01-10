package commands

import (
	"io"
	"os"

	"github.com/docker/machine/libmachine"
)

func cmdVersion(c CommandLine, api libmachine.API) error {
	return printVersion(c, api, os.Stdout)
}

func printVersion(c CommandLine, api libmachine.API, out io.Writer) error {
	if len(c.Args()) != 0 {
		return ErrTooManyArguments
	}

	c.ShowVersion()
	return nil
}
