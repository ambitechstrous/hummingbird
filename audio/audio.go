package audio

import (
	"fmt"
	"log"
	"math"

	"github.com/ambitechstrous/hummingbird/midi"
	"github.com/gordonklaus/portaudio"

	aubio "github.com/coral/aubio-go"
)

const sampleRate = 16000

type AudioProcessor struct {
	handler midi.IMidiHandler
	pitch   *aubio.Pitch
	stream  *portaudio.Stream
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

	// Aubio supports MIDI direct, but we want Hz to do conversions on voice harmonics
	pitch := aubio.NewPitch(aubio.PitchDefault, 2048, 1024, sampleRate)
	pitch.SetUnit(aubio.PitchOutFreq)
	pitch.SetTolerance(0.85)

	// Initialize processor with the provided MIDI handler
	processor := &AudioProcessor{
		handler: handler,
		pitch:   pitch,
	}

	// Create stream with mono input/output. Callback function will be called for each buffer of audio data.
	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, 1024, processor.processAudioInput)
	if err != nil {
		return nil, fmt.Errorf("Error opening PortAudio stream: %v", err)
	}

	// Stream set here instead of constructor, in order to utilize reference to processor in the callback function
	processor.stream = stream

	return processor, nil
}

// processAudioInput is the callback function for processing audio input. It is called by PortAudio and sends MIDI output to the handler.
func (p *AudioProcessor) processAudioInput(in []float32, out []float32) {
	// Convert samples to float64 for processing
	var sum float64
	samples := make([]float64, len(in))
	for i, v := range in {
		samples[i] = float64(v)
		sum += float64(v) * float64(v)
	}

	// RMS Gate. Avoid low energy buffers (i.e. could be bacvkground noise)
	rms := math.Sqrt(sum / float64(len(samples)))
	if rms < 0.001 {
		if err := p.handler.SendNote(0); err != nil {
			log.Printf("Error sending MIDI note: %v", err)
		}
		return
	}

	// Create a SimpleBuffer from the input samples for pitch detection
	buf := aubio.NewSimpleBufferData(uint(len(in)), samples)
	defer buf.Free()

	// Perform pitch detection on the buffer
	p.pitch.Do(buf)

	// Find the main frequency from the buffered samples
	freq := p.pitch.Buffer().Slice()[0]

	// Figure out the primary note, as the frequency includes harmonics
	base, octave := findPrimaryNote(freq)
	if base == 0.0 {
		return
	}
	correctedFreq := base * math.Pow(2, float64(octave))

	// Convert to MIDI note number
	noteNumber := 12*math.Log2(correctedFreq/440.0) + 69
	midiNum := uint8(math.Round(noteNumber))

	// Send MIDI note to the handler
	if err := p.handler.SendNote(midiNum); err != nil {
		log.Printf("Error sending MIDI note: %v", err)
	}
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
	p.pitch.Free() // Free the pitch object to avoid memory leaks

	// Close the audio stream to not allow any more audio processing
	if err := p.stream.Close(); err != nil {
		log.Printf("Error closing PortAudio stream: %v", err)
	}

	// Terminate PortAudio to clean up resources
	if err := portaudio.Terminate(); err != nil {
		log.Printf("Error terminating PortAudio: %v", err)
	}
}
