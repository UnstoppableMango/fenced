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
