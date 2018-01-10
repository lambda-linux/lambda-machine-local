extension = $(patsubst windows,.exe,$(filter windows,$(1)))

# Valid target combinations
VALID_OS_ARCH := "[darwin/amd64][linux/amd64][windows/amd64]"

os.darwin := darwin
os.linux := linux
os.windows := windows

arch.amd64 := amd64

define gocross
	$(if $(findstring [$(1)/$(2)],$(VALID_OS_ARCH)), \
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 \
		$(GO) build \
		-o $(PREFIX)/bin/lambda-machine-local_${os.$(1)}_${arch.$(2)}$(call extension,$(GOOS)) \
		-a $(VERBOSE_GO) -tags "static_build netgo $(BUILDTAGS)" -installsuffix netgo \
		-ldflags "$(GO_LDFLAGS) -extldflags -static" $(GO_GCFLAGS) ./cmd/machine.go;)
endef

build-clean:
	rm -Rf $(PREFIX)/bin/*

build-x: $(shell find . -type f -name '*.go')
	$(foreach GOARCH,$(TARGET_ARCH),$(foreach GOOS,$(TARGET_OS),$(call gocross,$(GOOS),$(GOARCH))))

$(PREFIX)/bin/lambda-machine-local$(call extension,$(GOOS)): $(shell find . -type f -name '*.go')
	$(GO) build \
	-o $@ \
	$(VERBOSE_GO) -tags "$(BUILDTAGS)" \
	-ldflags "$(GO_LDFLAGS)" $(GO_GCFLAGS) ./cmd/machine.go

build: $(PREFIX)/bin/lambda-machine-local$(call extension,$(GOOS))
