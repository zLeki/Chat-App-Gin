package helpers

import "github.com/microcosm-cc/bluemonday"

func Sanitize(s string) string {
	return bluemonday.NewPolicy().Sanitize(s) // This is a golang library that sanitizes user input
}
