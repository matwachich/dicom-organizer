package main

import (
	"fmt"
	"time"

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

	prefs := fyne.CurrentApp().Preferences()
	var timer *time.Timer

	w.entry = &widget.Entry{
		Text: prefs.String("destination"),
		OnChanged: func(s string) {
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(time.Second, func() {
				prefs.SetString("destination", w.entry.Text)
				fmt.Println("Destination saved")
			})
		},
		ActionItem: &widget.Button{Icon: theme.FolderOpenIcon(), Importance: widget.LowImportance, OnTapped: func() {
			if res, _ := zenity.SelectFile(zenity.Title("Choisier le dossier de destination"), zenity.Directory()); res != "" {
				w.entry.SetText(res)
			}
		}},
	}

	return w
}

func (w *folderSelect) Enable() {
	w.entry.Enable()
	w.entry.ActionItem.(*widget.Button).Enable()
}
func (w *folderSelect) Disable() {
	w.entry.Disable()
	w.entry.ActionItem.(*widget.Button).Disable()
}
func (w *folderSelect) Disabled() bool {
	return w.entry.Disabled()
}

func (w *folderSelect) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.entry)
}
