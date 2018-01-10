package commands

import (
	"github.com/blang/semver"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/log"
	"github.com/lambda-linux/lambda-machine-local/version"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const (
	slug = "lambda-linux/lambda-machine-local"
)

func cmdUpdateCli(c CommandLine, api libmachine.API) error {
	if len(c.Args()) != 0 {
		return ErrTooManyArguments
	}

	if log.GetDebug() {
		log.Debug("Enabling selfupdate.EnableLog()")
		selfupdate.EnableLog()
	}

	latest_gh, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		log.Info(err)
		return nil
	}
	if !found {
		log.Info("No release was found")
		return nil
	}

	current_ver := semver.MustParse(version.Version)

	if current_ver.GTE(latest_gh.Version) {
		log.Infof("Lambda Machine Local version %s is currently the latest version", current_ver)
		return nil
	}

	log.Info("Starting update...")
	latest_gh, err = selfupdate.UpdateSelf(current_ver, slug)
	if err != nil {
		log.Info(err)
		return nil
	}
	log.Infof("Updated to latest Lambda Machine Local version %s", latest_gh.Version)

	return nil
}
