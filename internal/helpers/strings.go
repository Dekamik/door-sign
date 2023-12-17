package helpers

import "unicode"

func Capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
