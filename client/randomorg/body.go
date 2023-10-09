package randomorg

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type BodyParser struct {
}

func NewBodyParser() BodyParser {
	return BodyParser{}
}

func (p BodyParser) ParseIntegers(bb []byte, contentType string) ([]int, error) {
	if !validContentType(contentType) {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	var integers []int
	bytesReader := bytes.NewReader(bb)
	bufReader := bufio.NewScanner(bytesReader)
	for bufReader.Scan() {
		if err := bufReader.Err(); err != nil {
			return nil, fmt.Errorf("failed to parse integers: %v", err)
		}
		line := bufReader.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		integer, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("failed to convert line to int: %s: %v", line, err)
		}
		integers = append(integers, integer)
	}

	return integers, nil
}

func validContentType(contentType string) bool {
	if contentType == "" {
		return false
	}
	elems := strings.Split(contentType, ";")
	return elems[0] == "text/plain"
}
