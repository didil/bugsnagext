package bugsnagext

import (
	"regexp"
	"runtime"

	"github.com/bugsnag/bugsnag-go"
)

const stackTracesBufferSize = 1 << 20 // 1MB

// RegisterGoroutinesInfoCallBack registers a callback to add goroutine stack traces to events
func RegisterGoroutinesInfoCallBack() {
	bugsnag.OnBeforeNotify(func(e *bugsnag.Event, c *bugsnag.Configuration) error {
		// init buffer for the stack traces
		buf := make([]byte, stackTracesBufferSize)
		// print stack traces for all go routines
		n := runtime.Stack(buf, true)
		// remove current goroutine stack trace
		stackTrace := cleanupStackTrace(string(buf[:n]))

		if stackTrace != "" {
			// add goroutines stacktraces if any exist
			e.MetaData.Add("goroutines", "stacktraces", stackTrace)
		}

		return nil
	})
}

var cleanupRegx = regexp.MustCompile(`goroutine \d+ \[\w+\]:`)

func cleanupStackTrace(stackTrace string) string {
	res := cleanupRegx.FindAllStringIndex(stackTrace, 2)

	// only 1 goroutine running
	if len(res) < 2 {
		return ""
	}

	// skip first (current) goroutine stack trace
	return stackTrace[res[1][0]:]
}
