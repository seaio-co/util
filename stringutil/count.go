package stringutil

import (
	"unicode"
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
		case isAlphabet(r):
			if !inWord {
				inWord = true
				n++
			}
		case inWord && (r == '\'' || r == '-'):
		default:
			inWord = false
		}

		str = str[size:]
	}
	return n
}

const minCJKCharacter = '\u3400'

// Checks r is a letter but not CJK character.
func isAlphabet(r rune) bool {
	if !unicode.IsLetter(r) {
		return false
	}

	switch {
	case r < minCJKCharacter:
		return true
	case r >= '\u4E00' && r <= '\u9FCC':
		return false
	case r >= '\u3400' && r <= '\u4D85':
		return false
	case r >= '\U00020000' && r <= '\U0002B81D':
		return false
	}
	return true
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
