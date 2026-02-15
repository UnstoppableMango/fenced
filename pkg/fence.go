package fenced

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

var (
	tildes    = []byte("~~~")
	backticks = []byte("```")
)

type Block struct {
	Content string
	Lang    string
}

func (b Block) String() string {
	return b.Content
}

func Parse(r io.Reader) ([]Block, error) {
	var blocks []Block
	var cur strings.Builder
	var lang string

	scanner := bufio.NewScanner(r)
	inBlock := false

	for scanner.Scan() {
		data := scanner.Bytes()
		if after, ok := cutPrefix(data); ok {
			if inBlock {
				blocks = append(blocks, Block{
					Content: cur.String(),
					Lang:    string(lang),
				})
				cur.Reset()
			}

			lang = string(after)
			inBlock = !inBlock
		} else if inBlock {
			cur.Write(data)
			cur.WriteRune('\n')
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

func cutPrefix(line []byte) (after []byte, ok bool) {
	if after, ok = bytes.CutPrefix(line, backticks); ok {
		return
	}
	return bytes.CutPrefix(line, tildes)
}
