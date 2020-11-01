package impl

import (
	"fmt"
	"strings"
)

func getSliceString(list []string) string {
	strBuilder := strings.Builder{}
	for _, str := range list {
		strBuilder.WriteString(fmt.Sprintf("\"%s\"", str))
	}
	return strBuilder.String()
}
