package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type DICOMOrganizerApp2 struct {
	widget.BaseWidget

	busy Busy

	dest *folderSelect
	tags *structureEntry

	removeSrc, overwriteDst *widget.Check
}

func newDICOMOrganizerApp2(win fyne.Window) *DICOMOrganizerApp2 {
	w := &DICOMOrganizerApp2{}
	w.ExtendBaseWidget(w)

	prefs := fyne.CurrentApp().Preferences()

	w.busy.Win = win

	return w
}

func (w *DICOMOrganizerApp2) CreateRenderer() fyne.WidgetRenderer {
	cfg := container.New(layout.NewFormLayout(),
		widget.NewLabel("Destination"), w.dest,
		widget.NewLabel("Structure"), w.tags,
		layout.NewSpacer(), container.NewHBox(w.removeSrc, w.overwriteDst),
	)

	return widget.NewSimpleRenderer(cfg)
}
