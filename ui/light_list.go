package ui

import (
	"fynecv/appdata"
	"fynecv/entity"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type LightWidget struct {
}

func NewLightWidget(data *appdata.AppData) fyne.CanvasObject {

	bound := binding.NewUntypedList()
	for _, l := range data.Lights {
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
