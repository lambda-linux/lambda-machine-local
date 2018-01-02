package main

import (
	"fmt"
	"os"
	"strconv"

	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/docker/machine/libmachine/drivers/plugin/localbinary"
	"github.com/docker/machine/libmachine/log"
	"github.com/lambda-linux/lambda-machine-local/commands"
	"github.com/lambda-linux/lambda-machine-local/commands/mcndirs"
	"github.com/lambda-linux/lambda-machine-local/drivers/virtualbox"
	"github.com/lambda-linux/lambda-machine-local/version"
)

var AppHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

{{.Usage}}

Version: {{.Version}}{{if or .Author .Email}}

Author:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`

var CommandHelpTemplate = `Usage: lambda-machine-local {{.Name}}{{if .Flags}} [OPTIONS]{{end}} [arg...]

{{.Usage}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .Flags}}

Options:
   {{range .Flags}}
   {{.}}{{end}}{{ end }}
`

func setDebugOutputLevel() {
	// TODO: I'm not really a fan of this method and really would rather
	// use -v / --verbose TBQH
	for _, f := range os.Args {
		if f == "-D" || f == "--debug" || f == "-debug" {
			log.SetDebug(true)
		}
	}

	debugEnv := os.Getenv("MACHINE_DEBUG")
	if debugEnv != "" {
		showDebug, err := strconv.ParseBool(debugEnv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing boolean value from MACHINE_DEBUG: %s\n", err)
			os.Exit(1)
		}
		log.SetDebug(showDebug)
	}
}

func main() {
	if os.Getenv(localbinary.PluginEnvKey) == localbinary.PluginEnvVal {
		driverName := os.Getenv(localbinary.PluginEnvDriverName)
		runDriver(driverName)
		return
	}

	localbinary.CurrentBinaryIsDockerMachine = true

	setDebugOutputLevel()
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = "Lambda Machine Local Contributors"
	app.Email = "https://github.com/lambda-linux/lambda-machine-local"

	app.Commands = commands.Commands
	app.CommandNotFound = cmdNotFound
	app.Usage = "Create and manage local machines running Docker."
	app.Version = version.FullVersion()

	log.Debug("Lamba Machine Local Version: ", app.Version)

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Enable debug mode",
		},
		cli.StringFlag{
			EnvVar: "MACHINE_STORAGE_PATH",
			Name:   "storage-path, s",
			Value:  mcndirs.GetBaseDir(),
			Usage:  "Configures storage path",
		},
		cli.StringFlag{
			EnvVar: "MACHINE_GITHUB_API_TOKEN",
			Name:   "github-api-token",
			Usage:  "Token to use for requests to the Github API",
			Value:  "",
		},
		cli.BoolFlag{
			EnvVar: "MACHINE_NATIVE_SSH",
			Name:   "native-ssh",
			Usage:  "Use the native (Go-based) SSH implementation.",
		},
		cli.StringFlag{
			EnvVar: "MACHINE_BUGSNAG_API_TOKEN",
			Name:   "bugsnag-api-token",
			Usage:  "BugSnag API token for crash reporting",
			Value:  "",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

func runDriver(driverName string) {
	switch driverName {
	case "virtualbox":
		plugin.RegisterDriver(virtualbox.NewDriver("", ""))
	default:
		fmt.Fprintf(os.Stderr, "Unsupported driver: %s\n", driverName)
		os.Exit(1)
	}
}

func cmdNotFound(c *cli.Context, command string) {
	log.Errorf(
		"%s: '%s' is not a %s command. See '%s --help'.",
		c.App.Name,
		command,
		c.App.Name,
		os.Args[0],
	)
	os.Exit(1)
}
