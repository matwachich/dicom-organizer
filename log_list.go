package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
		func() fyne.CanvasObject { return newLogListItem() },
		func(id widget.ListItemID, co fyne.CanvasObject) {
			co.(*logListItem).update(&w.data[id])
			w.list.SetItemHeight(id, co.MinSize().Height)
		},
	)

	return w
}

func (w *LogList) Reset() {
	w.data = w.data[:0]
	w.Refresh()
}

func (w *LogList) Append(src, dst string, err error) {
	w.data = append(w.data, logEntry{Src: src, Dst: dst, Err: err})
	w.list.Refresh()
	w.list.ScrollToBottom()
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

func newLogListItem() *logListItem {
	w := &logListItem{}
	w.ExtendBaseWidget(w)

	w.rt = &widget.RichText{
		Segments: []widget.RichTextSegment{
			&widget.TextSegment{
				Style: widget.RichTextStyle{TextStyle: fyne.TextStyle{Bold: true}},
			},
			&widget.TextSegment{
				Style: widget.RichTextStyle{TextStyle: fyne.TextStyle{Italic: true}},
			},
			&widget.TextSegment{
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
				},
			},
		},
		Scroll:     container.ScrollNone,
		Truncation: fyne.TextTruncateOff,
		Wrapping:   fyne.TextWrapWord,
	}

	w.icon = &widget.Icon{}

	return w
}

func (w *logListItem) update(data *logEntry) {
	w.rt.Segments[0].(*widget.TextSegment).Text = data.Src
	w.rt.Segments[1].(*widget.TextSegment).Text = "> " + data.Dst

	if data.Err == nil {
		w.rt.Segments[2].(*widget.TextSegment).Text = ""
		w.icon.Resource = theme.ConfirmIcon()
	} else {
		w.rt.Segments[2].(*widget.TextSegment).Text = data.Err.Error()
		w.icon.Resource = theme.ErrorIcon()
	}

	w.Refresh()
}

func (w *logListItem) Tapped(_ *fyne.PointEvent) {}

func (w *logListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, w.icon, nil, w.rt))
}
