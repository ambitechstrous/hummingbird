package midi

import (
	"fmt"
	"log"
	"sync"

	gomidi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi"
)

var _ = gomidi.NoteOn

type MidiHandler struct {
	midiOut  rtmidi.MIDIOut
	mutex    *sync.Mutex
	lastNote uint8
}

type IMidiHandler interface {
	SendNote(note uint8) error
	Close()
}

func findPortByName(midiOut rtmidi.MIDIOut, name string) (int, error) {
	portCount, err := midiOut.PortCount()
	if err != nil {
		return -1, err
	} else if portCount == 0 {
		return -1, fmt.Errorf("no MIDI output ports available")
	}

	for i := 0; i < portCount; i++ {
		portName, err := midiOut.PortName(i)
		if err != nil {
			return -1, err
		} else if portName == name {
			return i, nil
		}
	}

	return -1, fmt.Errorf("port not found: %s", name)
}

func NewMidiHandler(useVirtualPort bool) (IMidiHandler, error) {
	midiOut, err := rtmidi.NewMIDIOutDefault()
	if err != nil {
		return nil, err
	}

	if useVirtualPort {
		if err = midiOut.OpenVirtualPort("Hummingbird"); err != nil {
			return nil, err
		}
	} else {
		// FIXME: Make port name configurable or infer via machine config
		targetPort, err := findPortByName(midiOut, "IAC Driver Bus 1")
		if err != nil {
			return nil, err
		}

		if err = midiOut.OpenPort(targetPort, "Hummingbird"); err != nil {
			return nil, err
		}
	}

	return &MidiHandler{midiOut: midiOut, lastNote: 0, mutex: &sync.Mutex{}}, nil
}

func (h *MidiHandler) SendNote(note uint8) error {
	// Lock mutex to ensure notes are processed sequentially
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Check if a note was played previously and send a Note Off message for it
	if h.lastNote != note {
		// Send Note Off if the last note was not silence (i.e. non-zero)
		if h.lastNote != 0 {
			noteOff := gomidi.NoteOff(0, h.lastNote)
			return h.midiOut.SendMessage(noteOff)
		}

		// So long as the note isn't silence (i.e. non-zero), send the new note
		if note != 0 {
			noteOn := gomidi.NoteOn(0, note, 100) // Channel 0, velocity 100
			return h.midiOut.SendMessage(noteOn)
		}

		// Update the last note played to the current note
		h.lastNote = note
	}

	return nil
}

func (h *MidiHandler) Close() {
	if err := h.midiOut.Close(); err != nil {
		log.Printf("Error closing MIDI output: %v", err)
	}
}
