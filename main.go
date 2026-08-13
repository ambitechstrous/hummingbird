package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

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

	fmt.Println("Enter MIDI note numbers 0-127 (0 = note off). Crtl+C to quit.")

	// TODO: Audio processing loop instead of stdin scanner. This is just for testing purposes.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		num, err := strconv.ParseUint(line, 10, 8)
		if err != nil {
			fmt.Printf("Invalid note %d", num)
			continue
		} else if err := handler.SendNote(uint8(num)); err != nil {
			log.Printf("Error sending note: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}

	// Wait for termination signal to exit gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
