package main

import "fyne.io/fyne/v2/app"

func main() {
	a := app.New()
	w := a.NewWindow("DICOM Organizer")

	dicomApp := newDICOMOrganizerApp(w)
	w.SetCloseIntercept(func() {
		dicomApp.saveConfig()
		w.Close()
	})

	w.SetContent(dicomApp)
	w.ShowAndRun()
}
