# Validate DCO on all history
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

fmt:
	@test -z "$$(gofmt -s -l . 2>&1 | grep -v vendor/)"

vet:
	@test -z "$$(go vet $(PKGS) 2>&1)"

lint:
	$(if $(GOLINT), , \
		$(error Please install golint: go get -u github.com/golang/lint/golint))
	@test -z "$$($(GOLINT) ./... 2>&1 | grep -v vendor/ | grep -v "cli/" | grep -v "should have comment")"
