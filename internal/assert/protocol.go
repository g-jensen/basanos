package assert

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseProtocol(r io.Reader) (expected, actual string, err error) {
	reader := bufio.NewReader(r)

	version, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	if strings.TrimSpace(version) != "basanos:1" {
		return "", "", fmt.Errorf("invalid version header")
	}

	expected, err = readLengthPrefixedContent(reader)
	if err != nil {
		return "", "", err
	}

	actual, err = readLengthPrefixedContent(reader)
	if err != nil {
		return "", "", err
	}

	return expected, actual, nil
}

func BuildProtocol(expected, actual string) string {
	return fmt.Sprintf("basanos:1\n%d\n%s%d\n%s",
		len(expected), expected,
		len(actual), actual)
}

func readLengthPrefixedContent(reader *bufio.Reader) (string, error) {
	lengthLine, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	length, err := strconv.Atoi(strings.TrimSpace(lengthLine))
	if err != nil {
		return "", err
	}

	content := make([]byte, length)
	_, err = io.ReadFull(reader, content)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
