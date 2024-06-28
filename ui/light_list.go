package ui

import (
	"fynecv/entity"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewLightList(items []string) fyne.CanvasObject {
	lights := entity.BuildLightList(items)

	bound := binding.NewUntypedList()
	for _, l := range lights {
		bound.Append(l)
	}
	list := widget.NewListWithData(
		bound,
		func() fyne.CanvasObject {
			content := widget.NewLabel("label")
			return content
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			u, _ := i.(binding.Untyped).Get()
			light, _ := u.(*entity.Light)
			label.SetText(light.Attributes.Name)
		},
	)
	return list
}
