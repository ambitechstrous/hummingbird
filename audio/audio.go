package audio

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/dsp/fourier"

	"github.com/ambitechstrous/hummingbird/midi"
	"github.com/gordonklaus/portaudio"
)

var _ = portaudio.Initialize
var _ = fourier.NewFFT

// TODO: Audio processing loop instead of stdin scanner. This is just for testing purposes.
func ProcessAudioInput(handler midi.IMidiHandler) {
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
}
