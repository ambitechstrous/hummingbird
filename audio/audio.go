package audio

import (
	"gonum.org/v1/gonum/dsp/fourier"

	"github.com/gordonklaus/portaudio"
)

var _ = portaudio.Initialize
var _ = fourier.NewFFT
