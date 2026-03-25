package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type statusLabel struct {
	widget.BaseWidget

	l *widget.Label

	done, err int
}

func newStatusLabel() *statusLabel {
	w := &statusLabel{l: &widget.Label{}}
	w.ExtendBaseWidget(w)
	return w
}

func (w *statusLabel) Reset() {
	w.done, w.err = 0, 0
	w.Refresh()
}

func (w *statusLabel) Done() {
	w.done++
	w.Refresh()
}

func (w *statusLabel) Err() {
	w.err++
	w.Refresh()
}

func (w *statusLabel) Refresh() {
	w.l.SetText(strconv.Itoa(w.done) + " OK | " + strconv.Itoa(w.err) + " échec(s)")
}

func (w *statusLabel) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.l)
}
