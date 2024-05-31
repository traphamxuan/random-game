package utils

import "regexp"

var (
	UserNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,64}$`)
)
