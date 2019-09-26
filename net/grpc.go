package net

import (
	"fmt"
	"strings"
)

// ServiceMethod
func ServiceMethod(m string) (string, string, error) {
	if len(m) == 0 {
		return "", "", fmt.Errorf("malformed method name: %q", m)
	}

	if m[0] == '/' {
		parts := strings.Split(m, "/")
		if len(parts) != 3 || len(parts[1]) == 0 || len(parts[2]) == 0 {
			return "", "", fmt.Errorf("malformed method name: %q", m)
		}
		service := strings.Split(parts[1], ".")
		return service[len(service)-1], parts[2], nil
	}

	parts := strings.Split(m, ".")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed method name: %q", m)
	}
	return parts[0], parts[1], nil
}
