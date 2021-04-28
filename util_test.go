package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestAddTrailingSlashIfNotExistsOnGivenPath(t *testing.T) {
	testCases := []struct {
		desc string
		path string
		want string
	}{
		{
			desc: "rootdir",
			path: "/",
			want: "/",
		},
		{
			desc: "subdir with trailing slash",
			path: "/subdir/",
			want: "/subdir/",
		},
		{
			desc: "subdir without trailing slash",
			path: "/subdir",
			want: "/subdir/",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := AddTrailingSlashIfNotExistsOnGivenPath(tC.path)

			assert.Equal(tC.want, got)
		})
	}
}
