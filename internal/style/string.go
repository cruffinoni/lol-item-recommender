package style

import (
	"strings"
)

func ToTitleCase(s string) string {
	var titleCase string
	capitalizeNext := true

	for _, c := range s {
		if c == ' ' {
			capitalizeNext = true
			continue
		}

		if capitalizeNext {
			titleCase += strings.ToUpper(string(c))
			capitalizeNext = false
		} else {
			titleCase += strings.ToLower(string(c))
		}
	}

	return titleCase
}
