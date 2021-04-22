package busylight_sync

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func GetStructName(s interface{}) string {
	completeStructname := fmt.Sprintf("%T", s)
	splittedStructName := strings.Split(completeStructname, ".")

	return splittedStructName[len(splittedStructName)-1]
}

func GetFuncName(function interface{}) string {
	funcPointer := reflect.ValueOf(function).Pointer()
	return runtime.FuncForPC(funcPointer).Name()
}
