package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func main() {
}

func unpack(x string) (string, error) {
	if x == "" {
		return "", nil
	}

	if unicode.IsDigit(rune(x[0])) {
		return "", fmt.Errorf("invalid string")
	}

	res := ""
	isEscape := false

	for i := 0; i < len(x); i++ {
		if x[i] == '\\' {
			if isEscape {
				res += string(x[i])
				isEscape = false
				continue
			}
			isEscape = true
			continue
		}

		if !isEscape && unicode.IsDigit(rune(x[i])) {
			dig, _ := strconv.Atoi(string(x[i]))
			res += strings.Repeat(string(x[i-1]), dig-1)
			continue
		}

		res += string(x[i])
		isEscape = false
	}
	return res, nil
}
