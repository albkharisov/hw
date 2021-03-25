package hw02unpackstring

import (
	"errors"
	"strings"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var b strings.Builder
	var pc rune;	// previous character

	for _, c := range str {
		if (pc == 0) {
			if unicode.IsDigit(c) {
				return "", ErrInvalidString
			}
			pc = c
			continue
		}

		if unicode.IsDigit(c) {
			i, _ := strconv.Atoi(string(c))
			b.WriteString(strings.Repeat(string(pc), i))
			pc = 0
		} else {
			b.WriteRune(pc)
			pc = c
		}
	}
	if pc != 0 {
		b.WriteRune(pc)
	}

	return b.String(), nil
}

