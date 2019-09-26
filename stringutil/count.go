package stringutil

import (
	"unicode/utf8"
)

// Len
func Len(str string) int {
	return utf8.RuneCountInString(str)
}

// WordCount
func WordCount(str string) int {
	var r rune
	var size, n int
	inWord := false
	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)
		switch {
		case inWord && (r == '\'' || r == '-'):
		default:
			inWord = false
		}

		str = str[size:]
	}
	return n
}

// Width
func Width(str string) int {
	var r rune
	var size, n int

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)
		n += RuneWidth(r)
		str = str[size:]
	}

	return n
}

// RuneWidth
func RuneWidth(r rune) int {
	switch {
	case r == utf8.RuneError || r < '\x20':
		return 0
	case '\x20' <= r && r < '\u2000':
		return 1
	case '\u2000' <= r && r < '\uFF61':
		return 2
	case '\uFF61' <= r && r < '\uFFA0':
		return 1
	case '\uFFA0' <= r:
		return 2
	}
	return 0
}
