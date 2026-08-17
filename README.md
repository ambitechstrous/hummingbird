# hummingbird

A vst plugin with quick voice-to-midi melody conversion.  Hum or whistle a melody, and hummingbird will be able to create a midi track for you to use in your favorite DAW.

## Requirements

- Go 1.21+
- [Homebrew](https://brew.sh)

Install system dependencies:

```bash
brew install aubio portaudio rtmidi
```

## Building

```bash
make build
```

## Running

```bash
make run
```

## DAW Setup

Hummingbird outputs MIDI via the IAC Driver on macOS. To receive MIDI in your DAW:

1. Open **Audio MIDI Setup** (found in `/Applications/Utilities`)
2. Open the **IAC Driver** and check **Device is online**
3. In your DAW, create a Software Instrument track and set its MIDI input to **IAC Driver Bus 1**
4. Run Hummingbird, then arm the track and sing or hum a melody

## Contributors

Alwin Joy

David Chen

Ruben Madera

Jose Vargas

Aidan Pelisson