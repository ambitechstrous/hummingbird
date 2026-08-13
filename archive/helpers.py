import numpy as np

''' DICT CONTAINING LOWEST FREQUENCIES FOR EACH NOTE (i.e. 'base frequency') '''
def generateKeys():
	frequencies = dict()
	frequencies["NA"] = 0.00
	frequencies["C"] = 16.35
	frequencies["C#"] = 17.32
	frequencies["D"] = 18.35
	frequencies["D#"] = 19.45
	frequencies["E"] = 20.60
	frequencies["F"] = 21.83
	frequencies["F#"] = 23.12
	frequencies["G"] = 24.50
	frequencies["G#"] = 25.96
	frequencies["A"] = 27.50
	frequencies["A#"] = 29.14
	frequencies["B"] = 30.87
	return frequencies

''' HELPER FUNCTION TO FIND CLOSEST REAL NOTE (i.e. 'A', 'E', etc) GIVEN NOTE FREQUENCY '''
def closestNote(note):
	# Frequencies less than 0 are invalid in the scope of this program
	if note <= 0:
		return 0, 0

	# Initial Variables for min/max algorithm
	minimum_diff = -1
	closest_note = "NA"
	closest_octave = 0
	frequencies = generateKeys()
	for key in frequencies:
		# NA key used to prevent division by 0 error
		if key == "NA":
			continue

		# Divide current note by frequency given key. Find the note that is closest to its rounded value (i.e. 'most divisible by a real note').
		roundedDiv = round(note/frequencies[key])
		realDiv = note/frequencies[key]
		diff = abs(realDiv-roundedDiv)

		# Octave must be whole number
		octave = np.log2(roundedDiv)
		if (diff < minimum_diff or minimum_diff == -1) and float(octave) == float(round(octave)):
			minimum_diff = diff
			closest_octave = octave
			closest_note = key

	# Return base frequency and octave
	return frequencies[closest_note], closest_octave