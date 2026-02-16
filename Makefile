GO        ?= go
GOMOD2NIX ?= $(GO) tool gomod2nix
GINKGO    ?= $(GO) tool ginkgo
NIX       ?= nix
WATCHEXEC ?= watchexec

GO_SRC ?= $(shell find . -name '*.go')

build: bin/fenced
deps tidy: go.sum gomod2nix.toml
container ctr: bin/image.tar.gz

check:
	$(NIX) flake check

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
	mkdir -p $(dir $@) && ln -s $(abspath $<)/bin/fenced $@
else
bin/fenced: ${GO_SRC}
	$(GO) build -o $@
endif

bin/image.tar.gz: ${GO_SRC}
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
