package bugsnagext

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func Test_cleanupStackTrace_multi_goroutines(t *testing.T) {
	stackTrace := `goroutine 21 [running]:
main.main.func1(0xc000070000, 0xc000076000, 0x0, 0x0)
	/mypath/exnotify/cmd/generror/main.go:18 +0x6f
github.com/bugsnag/bugsnag-go.(*middlewareStack).runBeforeFilter(0x14e2070, 0x130ee38, 0xc000070000, 0xc000076000, 0x0, 0x0)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/middleware.go:52 +0x69
github.com/bugsnag/bugsnag-go.(*middlewareStack).Run(0x14e2070, 0xc000070000, 0xc000076000, 0xc000049c00, 0xc000070000, 0xc000076000)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/middleware.go:33 +0xbd
github.com/bugsnag/bugsnag-go.(*Notifier).NotifySync(0x14d8540, 0x1352040, 0xc00006a040, 0xc00006a000, 0xc00000e060, 0x2, 0x2, 0xc00000e060, 0x1)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/notifier.go:74 +0x18d
github.com/bugsnag/bugsnag-go.(*Notifier).Notify(0x14d8540, 0x1352040, 0xc00006a040, 0xc00000e060, 0x2, 0x2, 0x1, 0x2)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/notifier.go:56 +0xc9
github.com/bugsnag/bugsnag-go.Recover(0xc000049f90, 0x1, 0x2)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/bugsnag.go:145 +0x26d
panic(0x12a5e20, 0xc000012020)
	/usr/local/Cellar/go/1.15.2/libexec/src/runtime/panic.go:969 +0x175
main.AddOne()
	/mypath/exnotify/cmd/generror/main.go:70 +0x13a
main.Run.func1(0x1357020, 0xc0000aa008)
	/mypath/exnotify/cmd/generror/main.go:43 +0x8a
created by main.Run
	/mypath/exnotify/cmd/generror/main.go:40 +0x4d

goroutine 1 [runnable]:
main.main()
	/mypath/exnotify/cmd/generror/main.go:33 +0x178
	
goroutine 22 [sleep]:
time.Sleep(0x3b9aca00)
	/usr/local/Cellar/go/1.15.2/libexec/src/runtime/time.go:188 +0xbf
main.SleepOne()
	/mypath/exnotify/cmd/generror/main.go:78 +0x2a
main.Run.func2(0x1357020, 0xc0000aa008)
	/mypath/exnotify/cmd/generror/main.go:50 +0x8a
created by main.Run
	/mypath/exnotify/cmd/generror/main.go:47 +0x79

goroutine 23 [sleep]:
time.Sleep(0x77359400)
	/usr/local/Cellar/go/1.15.2/libexec/src/runtime/time.go:188 +0xbf
main.SleepTwo()
	/mypath/exnotify/cmd/generror/main.go:83 +0x2a
main.Run.func3(0x1357020, 0xc0000aa008)
	/mypath/exnotify/cmd/generror/main.go:57 +0x8a
created by main.Run
	/mypath/exnotify/cmd/generror/main.go:54 +0xa5`

	cleanST := cleanupStackTrace(stackTrace)

	expectedCleanST := `goroutine 1 [runnable]:
main.main()
	/mypath/exnotify/cmd/generror/main.go:33 +0x178
	
goroutine 22 [sleep]:
time.Sleep(0x3b9aca00)
	/usr/local/Cellar/go/1.15.2/libexec/src/runtime/time.go:188 +0xbf
main.SleepOne()
	/mypath/exnotify/cmd/generror/main.go:78 +0x2a
main.Run.func2(0x1357020, 0xc0000aa008)
	/mypath/exnotify/cmd/generror/main.go:50 +0x8a
created by main.Run
	/mypath/exnotify/cmd/generror/main.go:47 +0x79

goroutine 23 [sleep]:
time.Sleep(0x77359400)
	/usr/local/Cellar/go/1.15.2/libexec/src/runtime/time.go:188 +0xbf
main.SleepTwo()
	/mypath/exnotify/cmd/generror/main.go:83 +0x2a
main.Run.func3(0x1357020, 0xc0000aa008)
	/mypath/exnotify/cmd/generror/main.go:57 +0x8a
created by main.Run
	/mypath/exnotify/cmd/generror/main.go:54 +0xa5`

	assert.Equal(t, expectedCleanST, cleanST)
}

func Test_cleanupStackTrace_single_goroutine(t *testing.T) {
	stackTrace := `goroutine 21 [running]:
main.main.func1(0xc000070000, 0xc000076000, 0x0, 0x0)
	/mypath/exnotify/cmd/generror/main.go:18 +0x6f
github.com/bugsnag/bugsnag-go.(*middlewareStack).runBeforeFilter(0x14e2070, 0x130ee38, 0xc000070000, 0xc000076000, 0x0, 0x0)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/middleware.go:52 +0x69
github.com/bugsnag/bugsnag-go.(*middlewareStack).Run(0x14e2070, 0xc000070000, 0xc000076000, 0xc000049c00, 0xc000070000, 0xc000076000)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/middleware.go:33 +0xbd
github.com/bugsnag/bugsnag-go.(*Notifier).NotifySync(0x14d8540, 0x1352040, 0xc00006a040, 0xc00006a000, 0xc00000e060, 0x2, 0x2, 0xc00000e060, 0x1)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/notifier.go:74 +0x18d
github.com/bugsnag/bugsnag-go.(*Notifier).Notify(0x14d8540, 0x1352040, 0xc00006a040, 0xc00000e060, 0x2, 0x2, 0x1, 0x2)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/notifier.go:56 +0xc9
github.com/bugsnag/bugsnag-go.Recover(0xc000049f90, 0x1, 0x2)
	/gopath/pkg/mod/github.com/bugsnag/bugsnag-go@v1.5.3/bugsnag.go:145 +0x26d
panic(0x12a5e20, 0xc000012020)
	/usr/local/Cellar/go/1.15.2/libexec/src/runtime/panic.go:969 +0x175
main.AddOne()
	/mypath/exnotify/cmd/generror/main.go:70 +0x13a
main.Run.func1(0x1357020, 0xc0000aa008)
	/mypath/exnotify/cmd/generror/main.go:43 +0x8a
created by main.Run
	/mypath/exnotify/cmd/generror/main.go:40 +0x4d`

	cleanST := cleanupStackTrace(stackTrace)

	expectedCleanST := ``

	assert.Equal(t, expectedCleanST, cleanST)
}

func Test_RegisterGoroutinesInfoCallBack_notify(t *testing.T) {
	// register callback
	RegisterGoroutinesInfoCallBack()

	// start testserver
	s, ch := testutil.Setup()

	// configure bugsnag
	bugsnag.Configure(bugsnag.Configuration{
		APIKey: testutil.TestAPIKey,
		Endpoints: bugsnag.Endpoints{
			Notify: s.URL,
		},
		ProjectPackages: []string{"main"},
	})

	// spin goroutine
	go func() { time.Sleep(3 * time.Second) }()

	var wg sync.WaitGroup
	var payload []byte

	wg.Add(1)
	go func() {
		payload = <-ch
		wg.Done()
	}()

	go func() {
		//trigger test error
		bugsnag.Notify(fmt.Errorf("test error"))
	}()

	// wait till notification is issued
	time.Sleep(500 * time.Millisecond)

	// close server
	s.Close()

	// wait for the payload to be received
	wg.Wait()

	// parse payload
	stacktracesRes := gjson.Get(string(payload), "events.0.metaData.goroutines.stacktraces")
	stacktraces := stacktracesRes.String()
	assert.True(t, strings.HasPrefix(stacktraces, "goroutine "))
}

func Test_RegisterGoroutinesInfoCallBack_panic(t *testing.T) {
	// register callback
	RegisterGoroutinesInfoCallBack()

	// start testserver
	s, ch := testutil.Setup()

	// configure bugsnag
	bugsnag.Configure(bugsnag.Configuration{
		APIKey: testutil.TestAPIKey,
		Endpoints: bugsnag.Endpoints{
			Notify: s.URL,
		},
		ProjectPackages: []string{"main"},
	})

	// spin goroutine
	go func() { time.Sleep(3 * time.Second) }()

	var wg sync.WaitGroup
	var payload []byte

	wg.Add(1)
	go func() {
		payload = <-ch
		wg.Done()
	}()

	go func() {
		ctx := context.TODO()
		defer bugsnag.Recover(ctx)
		panic(fmt.Errorf("test error"))
	}()

	// wait till notification is issued
	time.Sleep(500 * time.Millisecond)

	// close server
	s.Close()

	// wait for the payload to be received
	wg.Wait()

	// parse payload
	stacktracesRes := gjson.Get(string(payload), "events.0.metaData.goroutines.stacktraces")
	stacktraces := stacktracesRes.String()
	assert.True(t, strings.HasPrefix(stacktraces, "goroutine "))
}
