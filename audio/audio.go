package audio

import (
	"fmt"
	"log"
	"math"

	"gonum.org/v1/gonum/dsp/fourier"

	"github.com/ambitechstrous/hummingbird/midi"
	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

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
	stream, err := portaudio.OpenDefaultStream(1, 1, sampleRate, 1024, processor.processAudioInput)
	if err != nil {
		return nil, fmt.Errorf("Error opening PortAudio stream: %v", err)
	}

	// Stream set here instead of constructor, in order to utilize reference to processor in the callback function
	processor.stream = stream

	return processor, nil
}

// FIXME: FFT Kinda sucks. Look into coral/aubio-go for better pitch detection.
func (p *AudioProcessor) processAudioInput(in []float32, out []float32) {
	// FFT helps us determine the frequency of the input signal.
	fft := fourier.NewFFT(len(in))

	// Convert samples to float64 for FFT processing
	samples := make([]float64, len(in))
	rms := 0.0
	for i := 0; i < len(in); i++ {
		samples[i] = float64(in[i])
		rms += float64(in[i]) * float64(in[i])
	}

	// Calculate RMS (Root Mean Square) of the input signal to determine its amplitude
	rms = math.Sqrt(rms / float64(len(in)))
	if rms < 0.01 {
		// If the signal is too quiet, we can consider it as silence
		return
	}

	// Perform FFT on the input audio data
	fftResult := fft.Coefficients(nil, samples)

	// Figure out which bin had the highest magnitude (i.e. the dominant frequency)
	maxMag := 0.0
	maxBin := 0
	for i, c := range fftResult {
		mag := math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
		if mag > maxMag {
			maxMag = mag
			maxBin = i
		}
	}

	// Calculate frequency corresponding to the bin index
	freq := fft.Freq(maxBin) * sampleRate

	// Convert frequency to MIDI note number
	midiNote := 12*math.Log2(freq/440) + 69
	midiNum := uint8(math.Round(midiNote))

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
	// Close the audio stream to not allow any more audio processing
	if err := p.stream.Close(); err != nil {
		log.Printf("Error closing PortAudio stream: %v", err)
	}

	// Terminate PortAudio to clean up resources
	if err := portaudio.Terminate(); err != nil {
		log.Printf("Error terminating PortAudio: %v", err)
	}
}
