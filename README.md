# fenced

[![CI](https://img.shields.io/github/actions/workflow/status/UnstoppableMango/fenced/ci.yml?branch=main&label=CI)](https://github.com/UnstoppableMango/fenced/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/UnstoppableMango/fenced/branch/main/graph/badge.svg)](https://codecov.io/gh/UnstoppableMango/fenced)
[![Go Report Card](https://goreportcard.com/badge/github.com/UnstoppableMango/fenced)](https://goreportcard.com/report/github.com/UnstoppableMango/fenced)
[![GoDoc](https://pkg.go.dev/badge/github.com/UnstoppableMango/fenced)](https://pkg.go.dev/github.com/UnstoppableMango/fenced)
[![License](https://img.shields.io/github/license/UnstoppableMango/fenced)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/UnstoppableMango/fenced)](go.mod)

Parse code fences from text.

## Usage

### Install

```shell
go install github.com/unstoppablemango/fenced@latest
```

Or add as a tool in your `go.mod`:

```shell
go get -tool github.com/unstoppablemango/fenced
```

Or download a binary from [GitHub Releases](https://github.com/UnstoppableMango/fenced/releases):

```shell
# Linux/macOS
curl -L https://github.com/UnstoppableMango/fenced/releases/latest/download/fenced_$(uname -s)_$(uname -m).tar.gz | tar xz
mkdir -p ~/.local/bin && mv fenced ~/.local/bin/
```

### CLI

```shell
$ fenced testdata/markdown.md
import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
```

### Docker

```shell
docker run -v $(pwd):/data ghcr.io/unstoppablemango/fenced:latest /data/testdata/markdown.md
```

Or with stdin:

```shell
cat testdata/markdown.md | docker run -i ghcr.io/unstoppablemango/fenced:latest
```

### Library

```go
import (
    "fmt"
    "io/fs"
    "os"

    fenced "github.com/unstoppablemango/fenced/pkg"
)

func main() {
    f, err := os.Open("testdata/markdown.md")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    blocks, err := fenced.Parse(f)
    if err != nil {
        panic(err)
    }

    for _, b := range blocks {
        fmt.Println(b)
    }
}
```
