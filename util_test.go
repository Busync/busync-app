package busylight_sync

import (
	"fmt"
	"strings"
)

func GetStructName(s interface{}) string {
	completeStructname := fmt.Sprintf("%T", s)
	splittedStructName := strings.Split(completeStructname, ".")

	return splittedStructName[len(splittedStructName)-1]
}
