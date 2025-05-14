package parser

import "regexp"

var re = regexp.MustCompile(`^[A-Za-z0-9*/_]+$`)

// At the moment we have same requirements for keys and values
func IsValidString(s string) bool {
	return re.MatchString(s)
}
