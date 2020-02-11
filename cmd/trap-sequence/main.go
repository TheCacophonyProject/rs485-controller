package main

import (
	"log"
	"os"
	"time"

	"github.com/TheCacophonyProject/rs485-controller/trapController"
	"github.com/TheCacophonyProject/window"
)

var (
	version = "<not set>"
)

func main() {
	err := runMain()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func runMain() error {
	log.SetFlags(0) // Removes timestamp output
	log.Printf("running version: %s", version)

	log.Println("starting service")
	if err := startDbusService(); err != nil {
		return err
	}
	w, err := window.New(
		"21:00",
		"06:00",
		0,
		0)
	if err != nil {
		return err
	}
	for {
		wait := w.Until() + time.Second
		log.Printf("waiting %v until starting sequence", wait)
		time.Sleep(wait)
		if err := trapController.StartSequence(); err != nil {
			return err
		}
		wait = w.UntilEnd() + time.Second
		log.Printf("waiting %v until stopping sequence", wait)
		time.Sleep(wait)
		if err := trapController.StopSequence(); err != nil {
			return err
		}
	}

}
