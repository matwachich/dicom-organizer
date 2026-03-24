package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/suyashkumar/dicom/pkg/tag"
)

type structureEntry struct {
	widget.BaseWidget

	entry *widget.Entry
}

func newStructureEntry() *structureEntry {
	w := &structureEntry{}
	w.ExtendBaseWidget(w)

	prefs := fyne.CurrentApp().Preferences()
	var timer *time.Timer

	w.entry = &widget.Entry{Text: prefs.String("structure")}
	w.entry.OnChanged = func(s string) {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(time.Second, func() {
			prefs.SetString("structure", w.entry.Text)
			fmt.Println("Structure saved")
		})
	}

	// default used tags
	tags := []tag.Tag{
		tag.PatientName, tag.PatientBirthDate, tag.PatientSex, tag.PatientAge,
		{},
		tag.StudyInstanceUID, tag.StudyDate, tag.StudyDescription,
		{},
		tag.SeriesInstanceUID, tag.SeriesNumber, tag.Modality, tag.SeriesDescription,
		{},
		tag.SOPInstanceUID, tag.InstanceNumber,
	}

	w.entry.ActionItem = &widget.Button{Icon: theme.ContentAddIcon(), Importance: widget.LowImportance, OnTapped: func() {
		menu := &fyne.Menu{}
		for _, t := range tags {
			if t.Compare(tag.Tag{}) == 0 {
				menu.Items = append(menu.Items, fyne.NewMenuItemSeparator())
				continue
			}

			info := tag.MustFind(t)
			name := info.Name
			keyword := info.Keyword

			menu.Items = append(menu.Items, &fyne.MenuItem{Label: name, Action: func() {
				w.insert("{" + keyword + "}")
				fyne.CurrentApp().Driver().CanvasForObject(w).Focus(w.entry)
			}})
		}

		widget.ShowPopUpMenuAtRelativePosition(menu, fyne.CurrentApp().Driver().CanvasForObject(w), fyne.NewPos(0, w.MinSize().Height+theme.Padding()), w.entry.ActionItem)
	}}

	return w
}

func (w *structureEntry) Enable() {
	w.entry.Enable()
	w.entry.ActionItem.(*widget.Button).Enable()
}
func (w *structureEntry) Disable() {
	w.entry.Disable()
	w.entry.ActionItem.(*widget.Button).Disable()
}
func (w *structureEntry) Disabled() bool {
	return w.entry.Disabled()
}

func (w *structureEntry) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.entry)
}

func (w *structureEntry) insert(text string) {
	for _, r := range text {
		w.entry.TypedRune(r)
	}
}
