package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func GetSliceOfStructNames(givenSlice []interface{}) []string {
	sliceOfStructNames := make([]string, 0)
	for _, element := range givenSlice {
		structName := GetStructName(element)
		sliceOfStructNames = append(sliceOfStructNames, structName)
	}

	return sliceOfStructNames
}

func GetStructName(s interface{}) string {
	completeStructname := fmt.Sprintf("%T", s)
	splittedStructName := strings.Split(completeStructname, ".")

	return splittedStructName[len(splittedStructName)-1]
}

func GetFuncName(function interface{}) string {
	funcPointer := reflect.ValueOf(function).Pointer()
	return runtime.FuncForPC(funcPointer).Name()
}
