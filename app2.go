package main

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
)

type DICOMOrganizerApp2 struct {
	widget.BaseWidget

	//busy Busy

	dest *folderSelect
	tags *structureEntry

	removeSrc, overwriteDst *widget.Check

	goFile, goFolder *widget.Button

	status *widget.Label

	log *LogList
}

func newDICOMOrganizerApp2(win fyne.Window) *DICOMOrganizerApp2 {
	w := &DICOMOrganizerApp2{}
	w.ExtendBaseWidget(w)

	//w.busy.Win = win
	prefs := fyne.CurrentApp().Preferences()

	w.dest = newFolderSelect()
	w.tags = newStructureEntry()

	w.removeSrc = &widget.Check{
		Text:      "Supprimer les fichiers source après la copie",
		OnChanged: func(b bool) { prefs.SetBool("removesrc", b) },
		Checked:   prefs.Bool("removesrc"),
	}

	w.overwriteDst = &widget.Check{
		Text:      "Ecraser les fichiers dans la destination",
		OnChanged: func(b bool) { prefs.SetBool("overwritedst", b) },
		Checked:   prefs.Bool("overwritedst"),
	}

	w.goFile = &widget.Button{Text: "Fichier(s)", Icon: theme.FileIcon(), Importance: widget.HighImportance, OnTapped: func() {
		if files, err := zenity.SelectFileMultiple(zenity.Title("Ajouter des fichiers DICOM"), zenity.FileFilters{
			zenity.FileFilter{Name: "Fichiers DICOM", Patterns: []string{"*.dcm", "*.dicom", "*.dic"}},
			zenity.FileFilter{Name: "Tous les fichiers", Patterns: []string{"*"}},
		}, zenity.Filename(prefs.String("lastfolder"))); err == nil {
			if len(files) > 0 {
				prefs.SetString("lastfolder", filepath.Dir(files[0]))
			}

			for _, f := range files {
				w.processFile(f)
			}
		}
	}}

	w.goFolder = &widget.Button{Text: "Dossier", Icon: theme.FolderIcon(), Importance: widget.HighImportance, OnTapped: func() {
		if folder, err := zenity.SelectFile(zenity.Title("Ajouter des fichiers DICOM"), zenity.Directory(), zenity.Filename(prefs.String("lastfolder"))); err == nil {
			prefs.SetString("lastfolder", folder)

		}
	}}

	w.log = newLogList()

	return w
}

func (w *DICOMOrganizerApp2) processFile(file string) {

}

func (w *DICOMOrganizerApp2) Enable() {
	w.dest.Enable()
	w.tags.Enable()
	w.removeSrc.Enable()
	w.overwriteDst.Enable()
	w.goFile.Enable()
	w.goFolder.Enable()
}
func (w *DICOMOrganizerApp2) Disable() {
	w.dest.Disable()
	w.tags.Disable()
	w.removeSrc.Disable()
	w.overwriteDst.Disable()
	w.goFile.Disable()
	w.goFolder.Disable()
}
func (w *DICOMOrganizerApp2) Disabled() bool {
	return w.goFile.Disabled()
}

func (w *DICOMOrganizerApp2) CreateRenderer() fyne.WidgetRenderer {
	cfg := container.New(layout.NewFormLayout(),
		widget.NewLabel("Destination"), w.dest,
		widget.NewLabel("Structure"), w.tags,
		layout.NewSpacer(), container.NewHBox(w.removeSrc, w.overwriteDst),
	)

	return widget.NewSimpleRenderer(container.NewBorder(
		container.NewVBox(cfg, container.NewBorder(nil, nil, container.NewHBox(w.goFile, w.goFolder), nil, w.status)),
		nil, nil, nil,
		w.log,
	))
}
