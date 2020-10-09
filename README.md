## Bugsnag Ext
bugsnag-go library extension to add goroutines stacktraces to an error event 

[![Build Status](https://travis-ci.org/didil/bugsnagext.svg?branch=master)](https://travis-ci.org/didil/bugsnagext)

### Run tests
```
$ make test
```
### Sample program / demo
Build and run the demo errors generator
```
$ make build-generror
# test with a manually notified/handled error
$ BUGSNAG_API_KEY=<BUGSNAG_API_KEY> bin/generror -e notif
# test with a panic/unhandled error
$ BUGSNAG_API_KEY=<BUGSNAG_API_KEY> bin/generror -e panic
```
Now you can see the new GOROUTINES tab in the bugsnag dashboard !
![Alt text](demo-assets/goroutines-tab.png?raw=true "Goroutines tab")

### Usage

Run `bugsnagext.RegisterGoroutinesInfoCallBack()` before your `bugsnag.Configure()` call

```
import (
...
	"github.com/bugsnag/bugsnag-go"
	"github.com/didil/bugsnagext"
...
)

func main() {
    // register goroutines info callback
    bugsnagext.RegisterGoroutinesInfoCallBack()
    // configure bugsnag
    bugsnag.Configure(bugsnag.Configuration{
        ...
    }}
    ...
}

```

#### Note
Logging handled errors in goroutines is transparent to the application code. But in order to log goroutines info properly during a goroutine's panic, you currently have to add `defer bugsnag.Recover(ctx)` at the beginning of the goroutine, otherwise the info about the other running goroutines is not registered. I couldn't yet find a way to make this work while being completely transparent to the application code. example:
```
go func() {
	// recover and log eventual panic
	defer bugsnag.Recover(ctx)

	MyFunction()
}()
```

