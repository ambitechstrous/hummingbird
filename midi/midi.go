package midi

import (
	"sync"

	gomidi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

var _ = gomidi.NoteOn

type MidiHandler struct {
	midiOut  drivers.Out
	mutex    *sync.Mutex
	lastNote uint8
}

type IMidiHandler interface {
	SendNote(note uint8)
}

func NewMidiHandler() (IMidiHandler, error) {
	drv, err := rtmididrv.New()
	if err != nil {
		return nil, err
	}

	out, err := drv.OpenVirtualOut("Hummingbird")
	if err != nil {
		return nil, err
	}

	return &MidiHandler{midiOut: out, lastNote: 0, mutex: &sync.Mutex{}}, nil
}

func (h *MidiHandler) SendNote(note uint8) {
	// Lock mutex to ensure notes are processed sequentially
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Check if a note was played previously and send a Note Off message for it
	if h.lastNote != note {
		// Send Note Off if the last note was not silence (i.e. non-zero)
		if h.lastNote != 0 {
			noteOff := []byte{0, h.lastNote, 0} // Channel 0, velocity 0
			h.midiOut.Send(noteOff)
		}

		// So long as the note isn't silence (i.e. non-zero), send the new note
		if note != 0 {
			noteOn := []byte{0, note, 100} // Channel 0, velocity 100
			h.midiOut.Send(noteOn)
		}

		// Update the last note played to the current note
		h.lastNote = note
	}
}
