package stringutil

import (
	"path"
	"strings"
)

// Reversed
func Reversed(str string) string {
	sr := []rune{}
	s1 := []rune(str)
	lens := len(s1)

	for i:=lens-1;i>=0;i-- {
		sr = append(sr,s1[i])
	}
	return string(sr)
}

// SplitQualifiedName
func SplitQualifiedName(str string) (string, string) {
	parts := strings.Split(str, "/")
	if len(parts) < 2 {
		return "", str
	}
	return parts[0], parts[1]
}

// JoinQualifiedName
func JoinQualifiedName(namespace, name string) string {
	return path.Join(namespace, name)
}

// ShortenString
func ShortenString(str string, n int) string {
	if len(str) <= n {
		return str
	}
	return str[:n]
}
