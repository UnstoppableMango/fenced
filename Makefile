GO         ?= go
GOMOD2NIX  ?= gomod2nix
GOLANGCI   ?= golangci-lint
GORELEASER ?= goreleaser
GINKGO     ?= ginkgo
NIX        ?= nix
PODMAN     ?= podman
WATCHEXEC  ?= watchexec

GO_SRC ?= $(shell find . -name '*.go')

build: bin/fenced
deps tidy: go.sum gomod2nix.toml
container ctr docker: bin/image.tar.gz

load: bin/stream_image.sh
	$< | $(PODMAN) load

check:
	$(NIX) flake check

lint:
	$(GOLANGCI) run

test:
	$(GINKGO) -r

watch:
	$(WATCHEXEC) -e go -- $(MAKE) test

cover: coverprofile.out
	$(GO) tool cover -func=$<

update:
	$(NIX) flake update

.PHONY: dist
dist:
	$(GORELEASER) build --snapshot --clean

clean:
	find . -type f -name '*cover*' -delete

ifneq (${IN_NIX_SHELL},)
bin/fenced: result
	mkdir -p ${@D} && ln -s $(abspath $<)/bin/fenced $@
else
bin/fenced: ${GO_SRC}
	$(GO) build -o $@
endif

bin/image.tar.gz: bin/stream_image.sh
	mkdir -p ${@D} && $< >$@

bin/stream_image.sh: ${GO_SRC}
	$(NIX) build .#ctr --out-link $@

go.sum: go.mod ${GO_SRC}
	$(GO) mod tidy
	@touch $@

gomod2nix.toml: go.mod go.sum
	$(GOMOD2NIX)

result: ${GO_SRC}
	$(NIX) build

coverprofile.out: ${GO_SRC}
	$(GINKGO) -r -cover
