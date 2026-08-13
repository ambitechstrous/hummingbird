#!/usr/bin/env python3
from tkinter import Tk, Frame, ttk
from audiohandler import AudioIO

APP_NAME = "HummingBird"

class Application(Frame):
	def __init__(self, master = None):
		super(Application, self).__init__(master)
		self.master.title(APP_NAME)
		self.parent = master
		self.buttonState = False
		self.audio = AudioIO()
		self.initUI()

	def initUI(self):
		style = ttk.Style()
		style.configure("TButton", padding=6, relief="flat", background="#66cc66d")
		self.mainButton = ttk.Button(text="Push", command=self.changeButtonState)
		self.mainButton.pack()

	def changeButtonState(self):
		self.buttonState = not(self.buttonState)
		if self.buttonState == True:
			self.audio.start()
		else:
			self.audio.stop()

def Main():
	# Initialize Application
	root = Tk()
	app = Application(root)

	# Set window size and background color
	app.parent.geometry('300x200+100+100')
	app.parent.configure(background = 'blue')
	app.mainloop()

if __name__ == "__main__":
	Main()