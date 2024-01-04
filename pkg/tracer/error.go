package tracer

import (
	"fmt"
	"runtime"
)

// AddStackTrace adds a stack trace to the given error and returns it.
func TraceError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w\n%s", err, getStackTrace())
}

// getStackTrace retrieves the stack trace information.
func getStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	stackTrace := "Stack Trace:\n"
	for {
		frame, more := frames.Next()
		stackTrace += fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}

	return stackTrace
}
