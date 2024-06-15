package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var res strings.Builder
	var letter rune
	var isBackSlash bool
	var isDigit bool
	for i, c := range input {
		switch {
		case unicode.IsDigit(c):
			if i == 0 || isDigit {
				return "", ErrInvalidString
			}
			if isBackSlash {
				AddRune(letter, &res)
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
		case c == '\\':
			isDigit = false
			AddRune(letter, &res)
			if isBackSlash {
				letter = c
				isBackSlash = false
			} else {
				letter = 0
				isBackSlash = true
			}
		default:
			isDigit = false
			if isBackSlash {
				return "", ErrInvalidString
			}
			AddRune(letter, &res)
			letter = c
		}
	}
	AddRune(letter, &res)
	return res.String(), nil
}

func AddRune(letter rune, res *strings.Builder) {
	if letter != 0 {
		res.WriteRune(letter)
	}
}
