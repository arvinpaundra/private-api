package util

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(ctx context.Context, duration time.Duration, operations map[string]func(ctx context.Context) error) <-chan struct{} {
	wait := make(chan struct{}, 1)

	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
		<-s

		close(s)

		timeFunc := time.AfterFunc(duration, func() {
			log.Printf("timeout %.2f has been elapsed, force exit\n", duration.Seconds())
			os.Exit(0)
		})

		defer timeFunc.Stop()

		for key, fn := range operations {
			err := fn(ctx)
			if err != nil {
				log.Printf("%s: failed cleaning up: %v\n", key, err.Error())
				return
			}

			log.Printf("%s was shutdown gracefully\n", key)
		}

		close(wait)
	}()

	return wait
}
