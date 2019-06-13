package util

import (
	"strings"
)

func TrimPK(public string) string {
	public = between(public, "KEY-----", "-----END")
	return strings.ToLower(strings.Join(strings.Split(strings.TrimSpace(public), "\n"), ""))
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirst := pos + len(a)
	if posFirst >= posLast {
		return ""
	}
	return value[posFirst:posLast]
}
