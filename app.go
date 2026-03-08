package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

type DICOMOrganizerApp struct {
	widget.BaseWidget

	busy Busy

	addF, addD, refresh, menu *widget.Button
	status                    *widget.Label

	data []sourceListItemData
	list *widget.List

	dest *folderSelect
	tags *tagsSelector

	removeSrc, overwriteDst *widget.Check

	start *widget.Button
}

func newDICOMOrganizerApp(win fyne.Window) *DICOMOrganizerApp {
	w := &DICOMOrganizerApp{}
	w.ExtendBaseWidget(w)

	prefs := fyne.CurrentApp().Preferences()

	w.busy.Win = win

	w.addF = widget.NewButtonWithIcon("", theme.FileIcon(), func() {
		if files, err := zenity.SelectFileMultiple(zenity.Title("Ajouter des fichiers DICOM"), zenity.FileFilters{
			zenity.FileFilter{Name: "Fichiers DICOM", Patterns: []string{"*.dcm", "*.dicom", "*.dic"}},
			zenity.FileFilter{Name: "Tous les fichiers", Patterns: []string{"*"}},
		}, zenity.Filename(prefs.String("lastfolder"))); err == nil {
			prefs.SetString("lastfolder", filepath.Dir(files[0]))

			for _, f := range files {
				w.loadFile(f)
			}

			w.list.Refresh()
			w.updateStatusCount()
		}
	})

	w.addD = widget.NewButtonWithIcon("", theme.FolderIcon(), func() {
		if folder, err := zenity.SelectFile(zenity.Title("Ajouter des fichiers DICOM"), zenity.Directory(), zenity.Filename(prefs.String("lastfolder"))); err == nil {
			prefs.SetString("lastfolder", folder)

			go func() {
				count := 0
				filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
					if !d.IsDir() {
						w.loadFile(path)

						count++
						fyne.Do(func() { w.busy.Set(count, -1) })
					}
					return nil
				})

				fyne.Do(func() {
					w.busy.Set(0, 0)

					w.list.Refresh()
					w.updateStatusCount()
				})
			}()
		}
	})

	w.refresh = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), w.updateStructure)

	w.menu = widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		widget.ShowPopUpMenuAtRelativePosition(&fyne.Menu{
			Items: []*fyne.MenuItem{
				{Label: "Sélectionner tout", Action: func() {
					for i := 0; i < len(w.data); i++ {
						w.data[i].selected = true
					}
					w.list.Refresh()
				}},
				{Label: "Désélectionner tout", Action: func() {
					for i := 0; i < len(w.data); i++ {
						w.data[i].selected = false
					}
					w.list.Refresh()
				}},
				{Label: "Inverser la sélection", Action: func() {
					for i := 0; i < len(w.data); i++ {
						w.data[i].selected = !w.data[i].selected
					}
					w.list.Refresh()
				}},
				fyne.NewMenuItemSeparator(),
				{Label: "Retirer la sélection", Action: func() {
					data := make([]sourceListItemData, 0, len(w.data))
					for i := 0; i < len(w.data); i++ {
						if !w.data[i].selected {
							data = append(data, w.data[i])
						}
					}
					w.data = data
					w.list.Refresh()
					w.updateStatusCount()
				}},
				{Label: "Retirer les terminés", Action: w.removeDone},
				{Label: "Réinitialiser les erreurs", Action: func() {
					for i := 0; i < len(w.data); i++ {
						if w.data[i].err != nil {
							w.data[i].err = nil
							w.data[i].done = false
						}
					}
					w.list.Refresh()
					w.updateStatusCount()
				}},
				{Label: "Retirer les doublons", Action: w.removeDuplicates},
				{Label: "Vider la liste", Action: func() {
					w.data = w.data[:0]
					w.list.Refresh()
					w.updateStatusCount()
				}},
			},
		}, win.Canvas(), fyne.NewPos(0, w.menu.MinSize().Height+theme.Padding()), w.menu)
	})

	w.status = &widget.Label{}

	w.list = widget.NewList(
		func() int { return len(w.data) },
		func() fyne.CanvasObject { return newSourceListItem(w) },
		func(id widget.ListItemID, co fyne.CanvasObject) {
			co.(*sourceListItem).update(id)
			w.list.SetItemHeight(id, co.MinSize().Height)
		},
	)

	w.dest = newFolderSelect()
	w.tags = newTagsSelector()

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

	w.start = &widget.Button{Text: "Démarrer", Importance: widget.HighImportance, OnTapped: func() {
		if len(w.data) <= 0 {
			return
		}

		dialog.ShowConfirm("Confirmation", "Voulez-vous retirer les éventuels doublons avant de lancer le traitement ?", func(b bool) {
			if b {
				w.removeDuplicates()
			}

			go func() {
				for i := 0; i < len(w.data); i++ {
					w.data[i].updateDestination(filepath.Join(w.dest.entry.Text, w.tags.entry.Text))
					w.data[i].process(w.removeSrc.Checked, w.overwriteDst.Checked)

					fyne.Do(func() {
						w.list.RefreshItem(i)
						w.busy.Set(i, len(w.data))
					})
				}

				fyne.Do(func() {
					w.busy.Set(0, 0)

					dialog.ShowConfirm("Traitement terminé", "Retirer les fichiers traités avec succès ?", func(b bool) {
						if b {
							w.removeDone()
						}
					}, win)
				})
			}()
		}, win)
	}}

	w.loadConfig()
	return w
}

func (w *DICOMOrganizerApp) updateStatusCount() {
	w.status.SetText(fmt.Sprintf("%d fichier(s) listé(s)", len(w.data)))
}

func (w *DICOMOrganizerApp) removeDone() {
	data := make([]sourceListItemData, 0, len(w.data))
	for i := 0; i < len(w.data); i++ {
		if !w.data[i].done {
			data = append(data, w.data[i])
		}
	}
	w.data = data
	w.list.Refresh()
	w.updateStatusCount()
}

func (w *DICOMOrganizerApp) removeDuplicates() {
	go func() {
		keys := make(map[string]struct{}, len(w.data))
		data := make([]sourceListItemData, 0, len(w.data))
		for i := 0; i < len(w.data); i++ {
			if _, ok := keys[w.data[i].source]; !ok {
				keys[w.data[i].source] = struct{}{}
				data = append(data, w.data[i])
			}

			fyne.Do(func() {
				w.busy.Set(i, len(w.data))
			})
		}
		w.data = data

		fyne.Do(func() {
			w.list.Refresh()
			w.updateStatusCount()
			w.busy.Set(0, 0)
		})
	}()
}

func (w *DICOMOrganizerApp) saveConfig() {
	pref := fyne.CurrentApp().Preferences()
	pref.SetString("destination", w.dest.entry.Text)
	pref.SetString("structure", w.tags.entry.Text)
}

func (w *DICOMOrganizerApp) loadConfig() {
	pref := fyne.CurrentApp().Preferences()
	w.dest.entry.Text = pref.String("destination")
	w.tags.entry.Text = pref.String("structure")
}

func (w *DICOMOrganizerApp) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		container.NewVBox(
			container.NewHBox(
				&widget.Label{Text: "Ajouter sources:", TextStyle: fyne.TextStyle{Bold: true}},
				w.addF, w.addD, w.refresh, w.status,
				layout.NewSpacer(),
				w.menu,
			),
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			container.New(layout.NewFormLayout(),
				widget.NewLabel("Destination"), w.dest,
				widget.NewLabel("Structure"), w.tags,
			),
			widget.NewSeparator(),
			container.NewHBox(w.removeSrc, w.overwriteDst),
			w.start,
		),
		nil, nil,
		w.list,
	))
}

// ----------------------------------------------------------------------------

type sourceListItemData struct {
	source    string
	dicomInfo dicom.Dataset

	destination string

	done bool
	err  error

	selected bool
}

func (w *sourceListItemData) updateDestination(destPath string) {
	var repl []string
	for _, tagKeyword := range regexp.MustCompile(`\{.+?\}`).FindAllString(destPath, -1) {
		if t, err := tag.FindByKeyword(strings.Trim(tagKeyword, "{}")); err == nil {
			if elem, err := w.dicomInfo.FindElementByTagNested(t.Tag); err == nil {
				repl = append(repl, tagKeyword)
				repl = append(repl, strings.Trim(elem.Value.String(), "[]"))
			}
		}
	}

	if len(repl) <= 0 {
		w.destination = ""
		return
	}

	w.destination = strings.NewReplacer(repl...).Replace(destPath)

	if absPath, _ := filepath.Abs(w.destination); absPath != "" {
		w.destination = absPath
	}

	elems := strings.Split(w.destination, string(filepath.Separator))
	if strings.HasSuffix(elems[0], ":") {
		elems[0] += string(filepath.Separator)
	}

	for i := 0; i < len(elems); i++ {
		elems[i] = strings.TrimSpace(elems[i])
	}

	w.destination = filepath.Clean(filepath.Join(elems...))
}

func (w *sourceListItemData) process(removeSrc, overwriteDst bool) (success bool) {
	if w.done || w.err != nil || w.destination == "" {
		return
	}

	w.err = doCopy(w.source, w.destination, removeSrc, overwriteDst)
	w.done = w.err == nil
	return
}

func (w *DICOMOrganizerApp) loadFile(file string) (err error) {
	itemData := sourceListItemData{
		source: file,
	}
	if itemData.dicomInfo, err = dicom.ParseFile(file, nil, dicom.SkipPixelData()); err != nil {
		itemData.err = err
	} else {
		itemData.updateDestination(filepath.Join(w.dest.entry.Text, w.tags.entry.Text))
	}

	w.data = append(w.data, itemData)
	return
}

func (w *DICOMOrganizerApp) updateStructure() {
	go func() {
		for i := 0; i < len(w.data); i++ {
			w.data[i].updateDestination(filepath.Join(w.dest.entry.Text, w.tags.entry.Text))

			fyne.Do(func() {
				w.busy.Set(i+1, len(w.data))
			})
		}

		fyne.Do(func() {
			w.list.Refresh()
		})
	}()
}

// ----------------------------------------------------------------------------

type sourceListItem struct {
	widget.BaseWidget

	parent *DICOMOrganizerApp
	id     widget.ListItemID

	sel    *widget.Check
	text   *widget.RichText
	status *widget.Icon
}

func newSourceListItem(parent *DICOMOrganizerApp) *sourceListItem {
	w := &sourceListItem{parent: parent}
	w.ExtendBaseWidget(w)

	w.sel = widget.NewCheck("", func(b bool) {
		w.parent.data[w.id].selected = b
	})

	w.text = &widget.RichText{
		Segments: []widget.RichTextSegment{
			&widget.TextSegment{},
			&widget.TextSegment{},
		},
		Scroll:     container.ScrollNone,
		Truncation: fyne.TextTruncateOff,
		Wrapping:   fyne.TextWrapWord,
	}

	w.status = widget.NewIcon(nil)

	return w
}

var (
	SuccessIcon = theme.NewSuccessThemedResource(theme.ConfirmIcon())
	ErrorIcon   = theme.NewErrorThemedResource(theme.ErrorIcon())
)

func (w *sourceListItem) update(id widget.ListItemID) {
	w.id = id
	w.sel.Checked = w.parent.data[id].selected

	w.text.Segments[0].(*widget.TextSegment).Text = w.parent.data[id].source
	w.text.Segments[0].(*widget.TextSegment).Style.TextStyle.Bold = true

	seg1 := w.text.Segments[1].(*widget.TextSegment)
	if w.parent.data[id].err == nil {
		seg1.Text = "> " + w.parent.data[id].destination
		seg1.Style.ColorName = ""
	} else {
		seg1.Text = "! " + w.parent.data[id].err.Error()
		seg1.Style.ColorName = theme.ColorNameError
	}

	if w.parent.data[w.id].done {
		w.status.Resource = SuccessIcon
	} else {
		if w.parent.data[w.id].err == nil {
			w.status.Resource = theme.ContentRemoveIcon()
		} else {
			w.status.Resource = ErrorIcon
		}
	}

	w.Refresh()
}

func (w *sourceListItem) Tapped(pe *fyne.PointEvent) {
	w.sel.SetChecked(!w.sel.Checked)
}

func (w *sourceListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, w.sel, w.status, w.text))
}
