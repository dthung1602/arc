package main

import (
	"github.com/dthung1602/arc/pkg/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting arc")
	arc, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		log.Println("Stopping arc")
		arc.Stop()
	}()

	arc.Serve()

	log.Println("Arc process finished")
}
