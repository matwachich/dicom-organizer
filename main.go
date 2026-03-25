package main

import "fyne.io/fyne/v2/app"

func main() {
	a := app.New()
	w := a.NewWindow("DICOM Organizer")

	dicomApp := newDICOMOrganizerApp()
	w.SetCloseIntercept(func() {
		dicomApp.savePrefs()
		w.Close()
	})

	w.SetContent(dicomApp)
	w.ShowAndRun()
}
