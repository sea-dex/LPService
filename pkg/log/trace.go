package log

import (
	"fmt"
	"strings"

	"github.com/hermeznetwork/tracerr"
)

func sprintStackTrace(st []tracerr.Frame) string {
	builder := strings.Builder{}
	// Skip deepest frame because it belongs to the go runtime and we don't
	// care about it.
	if len(st) > 0 {
		st = st[:len(st)-1]
	}

	for _, f := range st {
		builder.WriteString(fmt.Sprintf("\n%s:%d %s()", f.Path, f.Line, f.Func))
	}

	builder.WriteString("\n")

	return builder.String()
}

// appendStackTraceMaybeArgs will append the stacktrace to the args.
func appendStackTraceMaybeArgs(args []interface{}) []interface{} {
	for i := range args {
		if err, ok := args[i].(error); ok {
			err = tracerr.Wrap(err)
			st := tracerr.StackTrace(err)

			return append(args, sprintStackTrace(st))
		}
	}

	return args
}
