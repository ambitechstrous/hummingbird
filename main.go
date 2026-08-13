package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ambitechstrous/hummingbird/midi"
)

func main() {
	// FIXME: Virtual port should be in config instead of hardcoded
	handler, err := midi.NewMidiHandler(false)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := handler.Close(); err != nil {
			log.Printf("Error closing MIDI handler: %v", err)
		}
	}()

	// Example usage: Send a note (e.g., middle C, which is MIDI note number 60)
	if err := handler.SendNote(60); err != nil {
		log.Printf("Error sending note: %v", err)
	}

	time.Sleep(2 * time.Second) // Hold the note for 2 seconds
	if err := handler.SendNote(0); err != nil {
		log.Printf("Error sending note off: %v", err)
	}

	log.Println("Press Ctrl+C to exit...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
