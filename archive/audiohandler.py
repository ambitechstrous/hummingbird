''' AUDIO LIBRARIES '''
import rtmidi
import aubio
import pyaudio

''' MATH LIBRARIES '''
import math
import numpy as np

''' OTHER (BUILT-IN) LIBRARIES '''
import threading
import datetime

''' PERSONAL LIBRARIES '''
from helpers import closestNote

''' THREAD OBJECT FOR REAL-TIME AUDIO PROCESSING FUNCTIONALITY '''
class AudioIO(threading.Thread):
	def __init__(self, bitrate = 16000, chunk_size = 1024):
		super(AudioIO, self).__init__()
		self._stop = threading.Event()
		self.BITRATE = bitrate
		self.CHUNK_SIZE = chunk_size

	def run(self):
		self.createVirtualMidi()
		self.performDetection()

	def stop(self):
		self._stop.set()

	def stopped(self):
		return self._stop.isSet()

	def createVirtualMidi(self):
		# Use RtMIDI Library to create Virtual MIDI Instrument
		self._midiout = rtmidi.MidiOut()
		available_ports = self._midiout.get_ports()

		if available_ports:
			self._midiout.open_port(0)
		else:
			self._midiout.open_virtual_port("My virtual output")

	def performDetection(self):
		# Create stream object for Microphone input (using PyAudio)
		p = pyaudio.PyAudio()
		stream = p.open(format = pyaudio.paFloat32, channels = 1, rate = self.BITRATE, input = True, frames_per_buffer = self.CHUNK_SIZE)

		# Create detector object for Frequency Detection (using Aubio)
		detector = aubio.pitch("default", self.CHUNK_SIZE * 2, self.CHUNK_SIZE, self.BITRATE)
		detector.set_unit("Hz")
		detector.set_silence(-40)

		midiData = []

		# Loop until thread is stopped
		while not self.stopped():
			# Read from PyAudio input Stream
			tone = stream.read(CHUNK_SIZE)
			samples = np.fromstring(tone, dtype = aubio.float_type)

			# Perform pitch detection on samples. Creates array of size 1 with Frequency in HZ
			pitch = detector(samples)[0]
			volume = np.sum(samples**2)/len(samples)

			# Get closest real note given sung frequency
			base, octave = closestNote(pitch)
			if base == 0.00:
				continue
			freq = base*(2**octave)

			# Convert note frequency to MIDI Number
			noteNumber = 12*np.log2(freq/440.0) + 69 if freq > 0 else 0
			midiNum = int(round(noteNumber))

			# Send message to Virtual MIDI Instrument if pitch is changing
			if len(midiData) == 0:
				note_on = [0x90, midiNum, 100]
				self._midiout.send_message(note_on)
				midiData.append(midiNum)
			elif len(midiData) > 0 and midiData[len(midiData)-1] != midiNum:
				note_off = [0x80, midiData[len(midiData)-1], 0]
				note_on = [0x90, midiNum, 100]
				self._midiout.send_message(note_off)
				self._midiout.send_message(note_on)
				midiData.append(midiNum)

		# One Final Note Off Message to terminate final note
		if len(midiData) > 0:
			note_off = [0x80, midiData[len(midiData)-1], 0]
			self._midiout.send_message(note_off)

		# Close streams and objects
		stream.stop_stream()
		stream.close()
		p.terminate()

		# Delete Virtual MIDI Instrument to disconnect it from DAW
		del self._midiout