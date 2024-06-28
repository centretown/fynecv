package ui

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewBar() *widget.Toolbar {
	tb := widget.NewToolbar()
	tb.Append(widget.NewToolbarAction(theme.ColorChromaticIcon(), func() {}))
	return tb
}
