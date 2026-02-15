package fenced

import (
	"bufio"
	"io"
	"strings"
)

func Parse(r io.Reader) ([]string, error) {
	var codeBlocks []string
	scanner := bufio.NewScanner(r)
	inCodeBlock := false
	var codeBlock strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				codeBlocks = append(codeBlocks, codeBlock.String())
				codeBlock.Reset()
			}
			inCodeBlock = !inCodeBlock
		} else if inCodeBlock {
			codeBlock.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return codeBlocks, nil
}
