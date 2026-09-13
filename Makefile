GO_SRC ?= $(shell find . -name '*.go')

build: bin/fenced
deps tidy: go.sum gomod2nix.toml
container ctr docker: bin/image.tar.gz

load: bin/stream_image.sh
	$< | podman load

check:
	nix flake check

lint:
	golangci-lint run

test:
	go tool ginkgo -r

watch:
	watchexec -e go -- $(MAKE) test

cover: coverprofile.out
	go tool cover -func=$<

update:
	nix flake update

.PHONY: dist
dist:
	goreleaser build --snapshot --clean

clean:
	find . -type f -name '*cover*' -delete

ifneq (${IN_NIX_SHELL},)
bin/fenced: result
	mkdir -p ${@D} && ln -s $(abspath $<)/bin/fenced $@
else
bin/fenced: ${GO_SRC}
	go build -o $@
endif

bin/image.tar.gz: bin/stream_image.sh
	mkdir -p ${@D} && $< >$@

bin/stream_image.sh: ${GO_SRC}
	mkdir -p ${@D} && nix build .#ctr --out-link $@

go.sum: go.mod ${GO_SRC}
	go mod tidy
	@touch $@

gomod2nix.toml: go.mod go.sum
	gomod2nix

result: ${GO_SRC}
	nix build

coverprofile.out: ${GO_SRC}
	ginkgo -r -cover
