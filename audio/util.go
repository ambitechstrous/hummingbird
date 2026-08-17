package audio

import (
	"fmt"
	"math"
)

// Global variable to list out all valid whole notes
var validNotes = []float64{
	16.35, // C
	17.32, // C#
	18.35, // D
	19.45, // D#
	20.60, // E
	21.83, // F
	23.12, // F#
	24.50, // G
	25.96, // G#
	27.50, // A
	29.14, // A#
	30.87, // B
}

func findPrimaryNote(frequency float64) (float64, int) {
	fmt.Printf("Detected frequency: %f\n", frequency)

	// Frequencies less than 0 are considered invalid
	if frequency <= 0 {
		return 0.0, 0
	}

	closest_note := 0.0
	closest_octave := 0
	minimum_diff := math.MaxFloat64

	// Loop through all possible notes and find the frequency that is closest to the detected frequency
	for _, value := range validNotes {
		dividedFreq := frequency / value
		roundedDividedFreq := math.Round(dividedFreq)
		diff := math.Abs(dividedFreq - roundedDividedFreq)

		// Find the octave using power of 2 rule (frequency doubles every octave)
		octave := math.Log2(roundedDividedFreq)

		// When dealing with harmonic frequencies, the primary frequency is the closest note that is divisible by the octave
		if diff < minimum_diff && math.Abs(octave-math.Round(octave)) < 1e-10 {
			minimum_diff = diff
			closest_note = value
			closest_octave = int(octave)
		}
	}

	fmt.Printf("Closest note: %f, Closest octave: %d\n", closest_note, closest_octave)
	return closest_note, closest_octave
}
