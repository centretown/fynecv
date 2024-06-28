package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewAccordianList(objs ...fyne.CanvasObject) fyne.CanvasObject {
	acc := widget.NewAccordion()
	for i, o := range objs {
		text := fmt.Sprint("Item", i)
		acc.Append(widget.NewAccordionItem(text, o))
	}
	return acc
}
