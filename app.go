package main

import (
	"io/fs"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
	"github.com/suyashkumar/dicom"
)

type DICOMOrganizerApp struct {
	widget.BaseWidget

	dest *folderSelect
	tags *structureEntry

	removeSrc, overwriteDst *widget.Check

	goFile, goFolder *widget.Button

	status *statusLabel

	log *LogList
}

func newDICOMOrganizerApp() *DICOMOrganizerApp {
	w := &DICOMOrganizerApp{}
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

			w.processBegin()

			var dst string
			for _, f := range files {
				dst, err = w.processFile(f)
				w.processStep(f, dst, err)
			}

			w.processEnd()
		}
	}}

	w.goFolder = &widget.Button{Text: "Dossier", Icon: theme.FolderIcon(), Importance: widget.HighImportance, OnTapped: func() {
		if folder, err := zenity.SelectFile(zenity.Title("Ajouter des fichiers DICOM"), zenity.Directory(), zenity.Filename(prefs.String("lastfolder"))); err == nil {
			prefs.SetString("lastfolder", folder)

			w.processBegin()

			go func() {
				var dst string
				filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
					if !d.IsDir() {
						dst, err = w.processFile(path)
						w.processStep(path, dst, err)
					}
					return nil
				})

				w.processEnd()
			}()
		}
	}}

	w.status = newStatusLabel()

	w.log = newLogList()

	return w
}

func (w *DICOMOrganizerApp) processFile(file string) (dst string, err error) {
	dicomData, err := dicom.ParseFile(file, nil, dicom.SkipPixelData())
	if err != nil {
		return "", err
	}

	dst = sanitizePath(filepath.Join(w.dest.entry.Text, w.tags.getStructureForDicomData(dicomData)))

	err = doCopy(file, dst, w.removeSrc.Checked, w.overwriteDst.Checked)
	return
}

func (w *DICOMOrganizerApp) processBegin() {
	fyne.Do(func() {
		w.Disable()

		w.log.Reset()
		w.status.Reset()
	})
}

func (w *DICOMOrganizerApp) processEnd() {
	fyne.Do(func() {
		w.Enable()
	})
}

func (w *DICOMOrganizerApp) processStep(src, dst string, err error) {
	fyne.Do(func() {
		w.log.Append(src, dst, err)
		if err == nil {
			w.status.Done()
		} else {
			w.status.Err()
		}
	})
}

func (w *DICOMOrganizerApp) savePrefs() {
	w.dest.savePrefs()
	w.tags.savePrefs()
}

func (w *DICOMOrganizerApp) Enable() {
	w.dest.Enable()
	w.tags.Enable()
	w.removeSrc.Enable()
	w.overwriteDst.Enable()
	w.goFile.Enable()
	w.goFolder.Enable()
}
func (w *DICOMOrganizerApp) Disable() {
	w.dest.Disable()
	w.tags.Disable()
	w.removeSrc.Disable()
	w.overwriteDst.Disable()
	w.goFile.Disable()
	w.goFolder.Disable()
}
func (w *DICOMOrganizerApp) Disabled() bool {
	return w.goFile.Disabled()
}

func (w *DICOMOrganizerApp) CreateRenderer() fyne.WidgetRenderer {
	cfg := container.New(layout.NewFormLayout(),
		widget.NewLabel("Destination"), w.dest,
		widget.NewLabel("Structure"), w.tags,
		layout.NewSpacer(), container.NewHBox(w.removeSrc, w.overwriteDst),
		widget.NewLabel("Copier"), container.NewHBox(w.goFile, w.goFolder, w.status),
	)

	return widget.NewSimpleRenderer(container.NewBorder(
		container.NewVBox(cfg, widget.NewSeparator()),
		nil, nil, nil,
		w.log,
	))
}
