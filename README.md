# fenced

Parse code fences from text.

## Usage

```shell
$ fenced testdata/markdown.md
import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
```

```go
import (
    "fs"

    fenced "github.com/unstoppablemango/fenced/pkg"
)

func main() {
    f, err := fs.Open("testdata/markdown.md")
    if err != nil {
        panic(err)
    }

    blocks, err := fenced.Parse(f)
    if err != nil {
        panic(err)
    }

    for _, b := range blocks {
        fmt.Println(b)
    }
}
```
