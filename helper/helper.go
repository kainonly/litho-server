package helper

import "strings"

func StringToSlice(value string, sep string) []string {
	if value == "" {
		return []string{}
	}
	return strings.Split(value, sep)
}
