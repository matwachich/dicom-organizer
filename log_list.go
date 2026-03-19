package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type LogList struct {
	widget.BaseWidget

	data []logEntry
	list *widget.List
}

type logEntry struct {
	Src string
	Dst string
	Err error
}

func newLogList() *LogList {
	w := &LogList{}
	w.ExtendBaseWidget(w)

	w.list = widget.NewList(
		func() int { return len(w.data) },
	)

	return w
}

func (w *LogList) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.list)
}

//

type logListItem struct {
	widget.BaseWidget

	rt   *widget.RichText
	icon *widget.Icon
}
