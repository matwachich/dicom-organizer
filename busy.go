package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Busy struct {
	Win fyne.Window

	dlg  *dialog.CustomDialog
	prg  *widget.ProgressBar
	prgI *widget.ProgressBarInfinite
	lbl  *widget.Label
}

func (w *Busy) Set(count, total int) {
	if w.dlg == nil {
		w.prg = widget.NewProgressBar()
		w.prg.Hide()
		w.prgI = widget.NewProgressBarInfinite()
		w.prgI.Hide()

		w.lbl = &widget.Label{Alignment: fyne.TextAlignCenter, Text: "-"}

		w.dlg = dialog.NewCustomWithoutButtons("Traitement ...", container.NewVBox(container.NewStack(w.prg, w.prgI), w.lbl), w.Win)
	}

	if total < 0 {
		w.prg.Hide()
		w.prgI.Show()

		w.lbl.SetText(strconv.Itoa(count) + " ...")
	} else {
		w.prg.Show()
		w.prgI.Hide()

		w.lbl.SetText(strconv.Itoa(count) + "/" + strconv.Itoa(total))
		w.prg.SetValue(float64(count) / float64(total))
	}

	if count == total {
		w.dlg.Hide()
	} else {
		w.dlg.Show()
	}
}
