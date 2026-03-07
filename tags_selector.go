package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/suyashkumar/dicom/pkg/tag"
)

type tagsSelector struct {
	widget.BaseWidget

	entry *widget.Entry
}

func newTagsSelector() *tagsSelector {
	w := &tagsSelector{}
	w.ExtendBaseWidget(w)

	w.entry = &widget.Entry{}

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

func (w *tagsSelector) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.entry)
}

func (w *tagsSelector) insert(text string) {
	for _, r := range text {
		w.entry.TypedRune(r)
	}
}
