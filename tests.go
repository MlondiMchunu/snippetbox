package mlo

import (
	"fmt"
)

func containsString(v string, s []string) bool {
	for i, vs := range s {
		if v == vs {
			return true
		}
	}
	return false
}

func containsInt(v int, s []int) bool {
	for i, vs := range s {
		if v == vs {
			return true
		}
	}
	return false
}
func main() {
	fmt.Println()
}

// generics
func contains[T comparable](v T, s []T) bool {
	for i := range s {
		if v == s[i] {
			return true
		}
	}
	return false
}
