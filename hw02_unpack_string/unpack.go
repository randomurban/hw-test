package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var res strings.Builder
	var letter rune
	var isBackSlash bool
	var isDigit bool
	for i, c := range input {
		if c >= '0' && c <= '9' {
			if i == 0 || isDigit {
				return "", ErrInvalidString
			}
			if isBackSlash {
				if letter != 0 {
					res.WriteRune(letter)
				}
				letter = c
				isDigit = false
				isBackSlash = false
			} else {
				isDigit = true
				num, _ := strconv.Atoi(string(c))
				if num != 0 {
					res.WriteString(strings.Repeat(string(letter), num))
				}
				letter = 0
			}
		} else if c == '\\' {
			isDigit = false
			if letter != 0 {
				res.WriteRune(letter)
			}
			if isBackSlash {
				letter = c
				isBackSlash = false
			} else {
				letter = 0
				isBackSlash = true
			}
		} else {
			isDigit = false
			if isBackSlash {
				return "", ErrInvalidString
			}
			if letter != 0 {
				res.WriteRune(letter)
			}
			letter = c
		}

	}
	if letter != 0 {
		res.WriteRune(letter)
	}
	return res.String(), nil
}
