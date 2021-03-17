package music

import (
	"regexp"
	"strings"
)

func isNumber(str string) bool {
	result := false
	reg := regexp.MustCompile("^[0-9]+$")
	tmp := strings.Join(reg.FindAllString(str, 1), "")
	if tmp != "" {
		result = true
		return result
	}
	return result
}
