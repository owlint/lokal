package services

import (
	"os"
	"os/signal"
)

func HandleTermination(cleanUp func() int) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		os.Exit(cleanUp())
	}()
}

func WaitForTermination(cleanUp func() int) {
	HandleTermination(cleanUp)
	// wait indefinitely
	select {}
}
