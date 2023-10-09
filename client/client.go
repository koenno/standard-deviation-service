package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrSendRequest = errors.New("failed to send request")
	ErrResponse    = errors.New("response failure")

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

// // go:generate mockery --name=Converter --case underscore --with-expecter

type Client struct {
}

func New() Client {
	return Client{}
}

func (c Client) Send(req *http.Request) ([]byte, string, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrSendRequest, err)
	}

	defer resp.Body.Close()
	payloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("%w: unable to read body: %v", ErrResponse, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("%w: status code %d; body %s", ErrResponse, resp.StatusCode, string(payloadBytes))
	}

	return payloadBytes, resp.Header.Get("content-type"), nil
}

func ParseIntegers(bb []byte, contentType string) ([]int, error) {
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
