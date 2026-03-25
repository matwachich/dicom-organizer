package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

type structureEntry struct {
	widget.BaseWidget

	prefs fyne.Preferences
	r     *regexp.Regexp

	entry *widget.Entry
}

func newStructureEntry() *structureEntry {
	w := &structureEntry{
		prefs: fyne.CurrentApp().Preferences(),
		r:     regexp.MustCompile(`\{.+?\}`),
	}
	w.ExtendBaseWidget(w)

	var timer *time.Timer

	w.entry = &widget.Entry{Text: w.prefs.String("structure")}
	w.entry.OnChanged = func(s string) {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(time.Second, w.savePrefs)
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

func (w *structureEntry) getStructureForDicomData(dicomData dicom.Dataset) string {
	var repl []string
	for _, tagKeyword := range w.r.FindAllString(w.entry.Text, -1) {
		if t, err := tag.FindByKeyword(strings.Trim(tagKeyword, "{}")); err == nil {
			if elem, err := dicomData.FindElementByTagNested(t.Tag); err == nil {
				repl = append(repl, tagKeyword)
				repl = append(repl, forbiddenChars.Replace(strings.Trim(elem.Value.String(), "[]")))
			}
		}
	}

	return strings.NewReplacer(repl...).Replace(w.entry.Text)
}

func (w *structureEntry) savePrefs() {
	w.prefs.SetString("structure", w.entry.Text)
	fmt.Println("Structure saved")
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
