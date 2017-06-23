package commands

import (
	"errors"
	"flag"
	"os"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/crashreport"
	"github.com/docker/machine/libmachine/mcnerror"
	"github.com/stretchr/testify/assert"
)

func TestConsolidateError(t *testing.T) {
	cases := []struct {
		inputErrs   []error
		expectedErr error
	}{
		{
			inputErrs: []error{
				errors.New("Couldn't remove host 'bar'"),
			},
			expectedErr: errors.New("Couldn't remove host 'bar'"),
		},
		{
			inputErrs: []error{
				errors.New("Couldn't remove host 'bar'"),
				errors.New("Couldn't remove host 'foo'"),
			},
			expectedErr: errors.New("Couldn't remove host 'bar'\nCouldn't remove host 'foo'"),
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expectedErr, consolidateErrs(c.inputErrs))
	}
}

type MockCrashReporter struct {
	sent bool
}

func (m *MockCrashReporter) Send(err crashreport.CrashError) error {
	m.sent = true
	return nil
}

func TestSendCrashReport(t *testing.T) {
	defer func(fnOsExit func(code int)) { osExit = fnOsExit }(osExit)
	osExit = func(code int) {}

	defer func(factory func(baseDir string, apiKey string) crashreport.CrashReporter) {
		crashreport.NewCrashReporter = factory
	}(crashreport.NewCrashReporter)

	tests := []struct {
		description string
		err         error
		sent        bool
	}{
		{
			description: "Should send crash error",
			err: crashreport.CrashError{
				Cause:      errors.New("BUG"),
				Command:    "command",
				Context:    "context",
				DriverName: "virtualbox",
			},
			sent: true,
		},
		{
			description: "Should not send standard error",
			err:         errors.New("BUG"),
			sent:        false,
		},
	}

	for _, test := range tests {
		mockCrashReporter := &MockCrashReporter{}
		crashreport.NewCrashReporter = func(baseDir string, apiKey string) crashreport.CrashReporter {
			return mockCrashReporter
		}

		command := func(commandLine CommandLine, api libmachine.API) error {
			return test.err
		}

		context := cli.NewContext(cli.NewApp(), &flag.FlagSet{}, nil)
		runCommand(command)(context)

		assert.Equal(t, test.sent, mockCrashReporter.sent, test.description)
	}
}

func TestReturnExitCode1onError(t *testing.T) {
	command := func(commandLine CommandLine, api libmachine.API) error {
		return errors.New("foo is not bar")
	}

	exitCode := checkErrorCodeForCommand(command)

	assert.Equal(t, 1, exitCode)
}

func TestReturnExitCode3onErrorDuringPreCreate(t *testing.T) {
	if os.Getenv("BUGSNAG_ENABLE") == "" {
		t.Skip("skipping test; $BUGSNAG_ENABLE not set")
	}

	command := func(commandLine CommandLine, api libmachine.API) error {
		return crashreport.CrashError{
			Cause: mcnerror.ErrDuringPreCreate{
				Cause: errors.New("foo is not bar"),
			},
		}
	}

	exitCode := checkErrorCodeForCommand(command)

	assert.Equal(t, 3, exitCode)
}

func checkErrorCodeForCommand(command func(commandLine CommandLine, api libmachine.API) error) int {
	var setExitCode int

	originalOSExit := osExit

	defer func() {
		osExit = originalOSExit
	}()

	osExit = func(code int) {
		setExitCode = code
	}

	context := cli.NewContext(cli.NewApp(), &flag.FlagSet{}, nil)
	runCommand(command)(context)

	return setExitCode
}
