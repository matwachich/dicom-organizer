package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
)

type folderSelect struct {
	widget.BaseWidget

	entry *widget.Entry
}

func newFolderSelect() *folderSelect {
	w := &folderSelect{}
	w.ExtendBaseWidget(w)

	w.entry = widget.NewEntry()
	w.entry.ActionItem = &widget.Button{Icon: theme.FolderOpenIcon(), Importance: widget.LowImportance, OnTapped: func() {
		if res, _ := zenity.SelectFile(zenity.Title("Choisier le dossier de destination"), zenity.Directory()); res != "" {
			w.entry.SetText(res)
		}
	}}

	return w
}

func (w *folderSelect) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.entry)
}
