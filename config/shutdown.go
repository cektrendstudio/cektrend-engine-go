package config

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

var GlobalShutdown = &OnStop{}

type (
	shutdownFuncs func(ctx context.Context) error
	OnStop        struct {
		Shutdown map[string]shutdownFuncs
	}
)

func NewOnStop() *OnStop {
	return &OnStop{
		Shutdown: map[string]shutdownFuncs{},
	}
}

func (stp *OnStop) RegisterGracefullyShutdown(name string, operation shutdownFuncs) {
	if stp.Shutdown == nil {
		stp.Shutdown = map[string]shutdownFuncs{}
	}

	stp.Shutdown[name] = operation
}

func (stp *OnStop) GracefullyShutdown(ctx context.Context, timeout time.Duration) <-chan struct{} {
	wait := make(chan struct{})

	if stp.Shutdown == nil {
		return wait
	}

	go func() {
		var wg sync.WaitGroup

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		for key, op := range stp.Shutdown {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
