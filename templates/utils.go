package templates

import (
	"strings"
)

func escapeParamString(s string) string {
	s = strings.Replace(s, "&", "&amp;", -1)
	s = strings.Replace(s, "\"", "&quot;", -1)
	s = strings.Replace(s, "'", "&apos;", -1)
	s = strings.Replace(s, "<", "&lt;", -1)
	s = strings.Replace(s, ">", "&gt;", -1)
	return s
}

func compactTemplate(s string) string {
	var result string
	for _, line := range strings.Split(s, "\n") {
		if l := strings.TrimSpace(line); len(l) > 0 {
			result += l
		}
	}
	return result
}
