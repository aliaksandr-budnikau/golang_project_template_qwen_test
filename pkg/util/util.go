// Package util contains small, dependency-free helpers that may be useful to
// other modules importing this project.
package util

import "strings"

// Pluralize returns s with a trailing "s" unless it already ends in "s".
// It is a trivial example; replace it with real utilities as needed.
func Pluralize(s string) string {
	if strings.HasSuffix(s, "s") {
		return s
	}
	return s + "s"
}
