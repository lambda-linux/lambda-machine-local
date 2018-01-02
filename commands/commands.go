package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/crashreport"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnerror"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/persist"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/lambda-linux/lambda-machine-local/commands/mcndirs"
)

const (
	defaultMachineName = "default"
)

var (
	ErrHostLoad           = errors.New("All specified hosts had errors loading their configuration")
	ErrNoDefault          = fmt.Errorf("Error: No machine name(s) specified and no %q machine exists", defaultMachineName)
	ErrNoMachineSpecified = errors.New("Error: Expected to get one or more machine names as arguments")
	ErrExpectedOneMachine = errors.New("Error: Expected one machine name as an argument")
	ErrTooManyArguments   = errors.New("Error: Too many arguments given")

	osExit = func(code int) { os.Exit(code) }
)

// CommandLine contains all the information passed to the commands on the command line.
type CommandLine interface {
	ShowHelp()

	ShowVersion()

	Application() *cli.App

	Args() cli.Args

	IsSet(name string) bool

	Bool(name string) bool

	Int(name string) int

	String(name string) string

	StringSlice(name string) []string

	GlobalString(name string) string

	FlagNames() (names []string)

	Generic(name string) interface{}
}

type contextCommandLine struct {
	*cli.Context
}

func (c *contextCommandLine) ShowHelp() {
	cli.ShowCommandHelp(c.Context, c.Command.Name)
}

func (c *contextCommandLine) ShowVersion() {
	cli.ShowVersion(c.Context)
}

func (c *contextCommandLine) Application() *cli.App {
	return c.App
}

// targetHost returns a specific host name if one is indicated by the first CLI
// arg, or the default host name if no host is specified.
func targetHost(c CommandLine, api libmachine.API) (string, error) {
	if len(c.Args()) == 0 {
		defaultExists, err := api.Exists(defaultMachineName)
		if err != nil {
			return "", fmt.Errorf("Error checking if host %q exists: %s", defaultMachineName, err)
		}

		if defaultExists {
			return defaultMachineName, nil
		}

		return "", ErrNoDefault
	}

	return c.Args()[0], nil
}

func runAction(actionName string, c CommandLine, api libmachine.API) error {
	var (
		hostsToLoad []string
	)

	// If user did not specify a machine name explicitly, use the 'default'
	// machine if it exists.  This allows short form commands such as
	// 'docker-machine stop' for convenience.
	if len(c.Args()) == 0 {
		target, err := targetHost(c, api)
		if err != nil {
			return err
		}

		hostsToLoad = []string{target}
	} else {
		hostsToLoad = c.Args()
	}

	hosts, hostsInError := persist.LoadHosts(api, hostsToLoad)

	if len(hostsInError) > 0 {
		errs := []error{}
		for _, err := range hostsInError {
			errs = append(errs, err)
		}
		return consolidateErrs(errs)
	}

	if len(hosts) == 0 {
		return ErrHostLoad
	}

	if errs := runActionForeachMachine(actionName, hosts); len(errs) > 0 {
		return consolidateErrs(errs)
	}

	for _, h := range hosts {
		if err := api.Save(h); err != nil {
			return fmt.Errorf("Error saving host to store: %s", err)
		}
	}

	return nil
}

func runCommand(command func(commandLine CommandLine, api libmachine.API) error) func(context *cli.Context) {
	return func(context *cli.Context) {
		api := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())
		defer api.Close()

		if context.GlobalBool("native-ssh") {
			api.SSHClientType = ssh.Native
		}
		api.GithubAPIToken = context.GlobalString("github-api-token")
		api.Filestore.Path = context.GlobalString("storage-path")

		// TODO (nathanleclaire): These should ultimately be accessed
		// through the libmachine client by the rest of the code and
		// not through their respective modules.  For now, however,
		// they are also being set the way that they originally were
		// set to preserve backwards compatibility.
		mcndirs.BaseDir = api.Filestore.Path
		mcnutils.GithubAPIToken = api.GithubAPIToken
		ssh.SetDefaultClient(api.SSHClientType)

		if err := command(&contextCommandLine{context}, api); err != nil {
			log.Error(err)

			if crashErr, ok := err.(crashreport.CrashError); ok {
				crashReporter := crashreport.NewCrashReporter(mcndirs.GetBaseDir(), context.GlobalString("bugsnag-api-token"))
				crashReporter.Send(crashErr)

				if _, ok := crashErr.Cause.(mcnerror.ErrDuringPreCreate); ok {
					osExit(3)
					return
				}
			}

			osExit(1)
			return
		}
	}
}

func confirmInput(msg string) (bool, error) {
	fmt.Printf("%s (y/n): ", msg)

	var resp string
	_, err := fmt.Scanln(&resp)
	if err != nil {
		return false, err
	}

	confirmed := strings.Index(strings.ToLower(resp), "y") == 0
	return confirmed, nil
}

var Commands = []cli.Command{}

func printIP(h *host.Host) func() error {
	return func() error {
		ip, err := h.Driver.GetIP()
		if err != nil {
			return fmt.Errorf("Error getting IP address: %s", err)
		}

		fmt.Println(ip)

		return nil
	}
}

// machineCommand maps the command name to the corresponding machine command.
// We run commands concurrently and communicate back an error if there was one.
func machineCommand(actionName string, host *host.Host, errorChan chan<- error) {
	// TODO: These actions should have their own type.
	commands := map[string](func() error){}

	log.Debugf("command=%s machine=%s", actionName, host.Name)

	errorChan <- commands[actionName]()
}

// runActionForeachMachine will run the command across multiple machines
func runActionForeachMachine(actionName string, machines []*host.Host) []error {
	var (
		numConcurrentActions = 0
		errorChan            = make(chan error)
		errs                 = []error{}
	)

	for _, machine := range machines {
		numConcurrentActions++
		go machineCommand(actionName, machine, errorChan)
	}

	// TODO: We should probably only do 5-10 of these
	// at a time, since otherwise cloud providers might
	// rate limit us.
	for i := 0; i < numConcurrentActions; i++ {
		if err := <-errorChan; err != nil {
			errs = append(errs, err)
		}
	}

	close(errorChan)

	return errs
}

func consolidateErrs(errs []error) error {
	finalErr := ""
	for _, err := range errs {
		finalErr = fmt.Sprintf("%s\n%s", finalErr, err)
	}

	return errors.New(strings.TrimSpace(finalErr))
}
