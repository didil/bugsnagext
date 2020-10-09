package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/bugsnag/bugsnag-go"
	"github.com/didil/bugsnagext"
)

// ErrType to differrentiate error types
type ErrType string

const panicErrType ErrType = "panic"
const notifErrType ErrType = "notif"

// this program generates errors in a goroutine (running the buggy AddOne function), while other goroutines are running
// it allows to test the multi-goroutines stack track functionality agains a live bugsnag endpoint
func main() {
	// register goroutines info callback
	bugsnagext.RegisterGoroutinesInfoCallBack()

	bugsnagAPIKey := os.Getenv("BUGSNAG_API_KEY")
	if bugsnagAPIKey == "" {
		fmt.Printf("please set the env variable BUGSNAG_API_KEY\n")
		os.Exit(1)
	}

	// configure bugsnag
	bugsnag.Configure(bugsnag.Configuration{
		APIKey: bugsnagAPIKey,
		// The import paths for the Go packages containing your source files
		ProjectPackages: []string{"main"},
	})

	// parse error type. notifErrType triggers a handled error while panicErrType triggers a panic
	errTypeStr := flag.String("e", string(notifErrType), fmt.Sprintf("error type: '%s' (handled error) or '%s' (unhandled)", notifErrType, panicErrType))
	flag.Parse()

	errType := ErrType(*errTypeStr)
	fmt.Printf("errType: %v\n", errType)

	if errType != notifErrType && errType != panicErrType {
		fmt.Printf("invalid errtype: %v\n", errType)
		os.Exit(1)
	}

	ctx := context.TODO()
	Run(ctx, errType)
}

// Run app
func Run(ctx context.Context, errType ErrType) {
	// launch AddOne
	go func() {
		// recover and log eventual panic
		defer bugsnag.Recover(ctx)

		AddOne(errType)
	}()

	wg := sync.WaitGroup{}
	// launch SleepOne
	wg.Add(1)
	go func() {
		SleepOne()
		wg.Done()
	}()

	// launch SleepTwo
	wg.Add(1)
	go func() {
		SleepTwo()
		wg.Done()
	}()
	wg.Wait()
}

// AddOne is a buggy function that triggers errors
func AddOne(errType ErrType) {
	rand.Seed(time.Now().UnixNano())
	limit := rand.Int31n(10000)

	for x := 0; ; x++ {
		if x >= int(limit) {
			switch errType {
			case panicErrType:
				panic(fmt.Errorf("[panic] x: %d", x))
			case notifErrType:
				bugsnag.Notify(fmt.Errorf("[notif] x: %d", x))
			}

			return
		}
	}
}

// SleepOne sleeps
func SleepOne() {
	time.Sleep(2 * time.Second)
	fmt.Printf("SleepOne done\n")
}

// SleepTwo is another sleeping function
func SleepTwo() {
	time.Sleep(3 * time.Second)
	fmt.Printf("SleepTwo done\n")
}
