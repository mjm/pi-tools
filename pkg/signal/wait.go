package signal

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Wait() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch
	log.Printf("Shutting down...")
}
