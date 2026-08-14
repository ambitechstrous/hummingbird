package audio

import (
	"fmt"
	"log"

	"gonum.org/v1/gonum/dsp/fourier"

	"github.com/ambitechstrous/hummingbird/midi"
	"github.com/gordonklaus/portaudio"
)

var _ = portaudio.Initialize
var _ = fourier.NewFFT

type AudioProcessor struct {
	stream  *portaudio.Stream
	handler midi.IMidiHandler
}

type IAudioProcessor interface {
	Start() error
	Stop() error
	Close()
}

func NewAudioProcessor(handler midi.IMidiHandler) (IAudioProcessor, error) {
	err := portaudio.Initialize()
	if err != nil {
		return nil, fmt.Errorf("Error initializing PortAudio: %v", err)
	}

	// Initialize processor with the provided MIDI handler
	processor := &AudioProcessor{
		handler: handler,
	}

	// Create stream with mono input/output, 44100 HZ sample rate. Callback function will be called for each buffer of audio data.
	stream, err := portaudio.OpenDefaultStream(1, 1, 44100, portaudio.FramesPerBufferUnspecified, processor.processAudioInput)
	if err != nil {
		return nil, fmt.Errorf("Error opening PortAudio stream: %v", err)
	}

	// Stream set here instead of constructor, in order to utilize reference to processor in the callback function
	processor.stream = stream

	return processor, nil
}

func (p *AudioProcessor) processAudioInput(in []float32, out []float32) {
	log.Printf("Received audio input: %v", in)
	log.Printf("Received output %v", out)
}

func (p *AudioProcessor) Start() error {
	if err := p.stream.Start(); err != nil {
		return fmt.Errorf("Error starting PortAudio stream: %v", err)
	}
	return nil
}

func (p *AudioProcessor) Stop() error {
	if err := p.stream.Stop(); err != nil {
		return fmt.Errorf("Error stopping PortAudio stream: %v", err)
	}
	return nil
}

func (p *AudioProcessor) Close() {
	// Close the audio stream to not allow any more audio processing
	if err := p.stream.Close(); err != nil {
		log.Printf("Error closing PortAudio stream: %v", err)
	}

	// Terminate PortAudio to clean up resources
	if err := portaudio.Terminate(); err != nil {
		log.Printf("Error terminating PortAudio: %v", err)
	}
}
