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

type AudioProcessor struct {
	handler midi.IMidiHandler
}

type IAudioProcessor interface {
	ProcessAudioInput()
	Close()
}

func NewAudioProcessor(handler midi.IMidiHandler) (IAudioProcessor, error) {
	err := portaudio.Initialize()
	if err != nil {
		return nil, fmt.Errorf("Error initializing PortAudio: %v", err)
	}

	return &AudioProcessor{
		handler: handler,
	}, nil
}

// TODO: Audio processing loop instead of stdin scanner. This is just for testing purposes.
func (p *AudioProcessor) ProcessAudioInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		num, err := strconv.ParseUint(line, 10, 8)
		if err != nil {
			fmt.Printf("Invalid note %d", num)
			continue
		} else if err := p.handler.SendNote(uint8(num)); err != nil {
			log.Printf("Error sending note: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}

func (p *AudioProcessor) Close() {
	if err := portaudio.Terminate(); err != nil {
		log.Printf("Error terminating PortAudio: %v", err)
	}
}
