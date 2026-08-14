package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ambitechstrous/hummingbird/audio"
	"github.com/ambitechstrous/hummingbird/midi"
)

func main() {
	// FIXME: Virtual port should be in config instead of hardcoded
	handler, err := midi.NewMidiHandler(false)
	if err != nil {
		panic(err)
	}

	defer handler.Close()

	processor, err := audio.NewAudioProcessor(handler)
	if err != nil {
		panic(err)
	}

	defer processor.Close()

	fmt.Println("Enter MIDI note numbers 0-127 (0 = note off). Crtl+C to quit.")
	if err := processor.Start(); err != nil {
		panic(err)
	}

	defer processor.Stop()

	// Wait for termination signal to exit gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
