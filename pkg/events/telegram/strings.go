package telegram

import "strings"

func concatStrings(strs ...string) string {
	b := strings.Builder{}
	for i, str := range strs {
		b.WriteString(str)
		if i < len(strs)-1 {
			b.WriteString(" - ")
		}
	}
	return b.String()
}
