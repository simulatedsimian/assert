package testbuddy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"strings"
)

func ConvertTo(i interface{}, to reflect.Type) (interface{}, error) {
	from := reflect.TypeOf(i)
	if from.ConvertibleTo(to) {
		return reflect.ValueOf(i).Convert(to), nil
	}
	return nil, fmt.Errorf("Cannot Convert From %v to %v", from, to)
}

func GetFullFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetShortFuncName(i interface{}) string {
	name := GetFullFuncName(i)
	return name[strings.LastIndex(name, ".")+1:]
}

type CallerInfo struct {
	pc       uintptr
	filename string
	lineNum  int
	lineSrc  string
}

func GetCallerInfo(skip int) (CallerInfo, bool) {
	pc, file, line, ok := runtime.Caller(skip + 1)

	ci := CallerInfo{pc, file, line, "[no source]"}
	if ok {
		data, err := ioutil.ReadFile(file)
		if err == nil {
			lines := bytes.Split(data, []byte{'\n'})
			ci.lineSrc = string(lines[ci.lineNum-1])
		}
	}
	return ci, ok
}

func SourceInfo(skip ...int) string {
	sk := 1
	if len(skip) > 0 {
		sk = skip[0]
	}
	ci, _ := GetCallerInfo(sk)
	return fmt.Sprintf("%s : [%d]\n%s", ci.filename, ci.lineNum, strings.Trim(ci.lineSrc, " \t"))
}