package main

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/hass"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type BindingData struct {
	Pan  binding.Float
	Tilt binding.Float
}

type View struct {
	Data      *appdata.AppData
	Binder    BindingData
	Container *fyne.Container
	Image     *canvas.Image
	Pan       *widget.Slider
	Tilt      *widget.Slider
}

func NewView(data *appdata.AppData) *View {

	mv := &View{
		Data: data,
		Binder: BindingData{
			Pan:  binding.NewFloat(),
			Tilt: binding.NewFloat(),
		},
		Image: canvas.NewImageFromImage(image.NewNRGBA(image.Rect(0, 0, 1280, 720))),
	}

	mv.Pan = widget.NewSliderWithData(0, 180, mv.Binder.Pan)
	mv.Tilt = widget.NewSliderWithData(0, 180, mv.Binder.Tilt)
	mv.Bind()

	mv.Image.FillMode = canvas.ImageFillContain
	mv.Tilt.Orientation = widget.Vertical

	mv.Container = container.NewBorder(
		nil,
		mv.Pan,
		nil,
		mv.Tilt,
		mv.Image)

	return mv
}

func (view *View) Bind() {
	view.Binder.Pan.Set(140)
	// view.Binder.Tilt.Set(50)

	view.Binder.Pan.AddListener(binding.NewDataListener(func() {
		value, _ := view.Binder.Pan.Get()
		cmd := fmt.Sprintf(`{"entity_id": "number.pan", "value": %.0f}`, value)
		hass.Post("services/number/set_value", cmd)
	}))
	view.Binder.Tilt.AddListener(binding.NewDataListener(func() {
		value, _ := view.Binder.Tilt.Get()
		cmd := fmt.Sprintf(`{"entity_id": "number.tilt", "value": %.0f}`, value)
		hass.Post("services/number/set_value", cmd)
	}))

}
