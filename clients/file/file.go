package file

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func (fileClient *Client) read(filepath string, resultsChannel chan<- string) (int, error) {
	count := 0
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return count, err
		}
		count++
		resultsChannel <- strings.TrimSpace(line)
		if err == io.EOF {
			break
		}
	}

	return count, nil
}
