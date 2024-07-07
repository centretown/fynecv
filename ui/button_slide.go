package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ButtonSlide struct {
	*fyne.Container
	Decrease *widget.Button
	Increase *widget.Button
	Slider   *widget.Slider
	Icon     *widget.Icon
}

func NewButtonSlide(field binding.Float, bounds *FloatEntryBounds) *ButtonSlide {
	buttonCheck := func(val float64) func() {
		return IncrementFloat(val, field, bounds)
	}

	bs := &ButtonSlide{
		Decrease: widget.NewButtonWithIcon("", theme.MoveDownIcon(), buttonCheck(-1)),
		Increase: widget.NewButtonWithIcon("", theme.MoveUpIcon(), buttonCheck(1)),
		Slider:   widget.NewSliderWithData(bounds.MinVal, bounds.MaxVal, field),
	}

	bs.Container = container.NewBorder(nil, nil, bs.Decrease, bs.Increase, bs.Slider)
	return bs
}
