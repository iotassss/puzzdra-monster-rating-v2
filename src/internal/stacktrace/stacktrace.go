package stacktrace

import (
	"fmt"
	"runtime"
	"strings"
)

func Print() string {
	var stackTrace []string

	const depth = 10
	for i := 1; i < depth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcName := runtime.FuncForPC(pc).Name()
		stackTrace = append(stackTrace, fmt.Sprintf("%s:%d %s", file, line, funcName))
	}

	return strings.Join(stackTrace, " ")
}
