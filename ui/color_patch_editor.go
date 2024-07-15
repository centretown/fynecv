package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ColorPatchEditor struct {
	Dialog   *dialog.CustomDialog
	window   fyne.Window
	onUpdate func()

	source      *ColorPatch
	patch       *ColorPatch
	applyButton *widget.Button

	hue          binding.Float
	saturation   binding.Float
	value        binding.Float
	unused       binding.Bool
	removeButton *widget.Button

	hsv HSV
}

func NewColorPatchEditor(source *ColorPatch,
	window fyne.Window, onUpdate func()) *ColorPatchEditor {

	pe := &ColorPatchEditor{
		source:   source,
		patch:    NewColorPatchWithColor(source.FillColor, func() {}, func() {}),
		window:   window,
		onUpdate: onUpdate,

		hue:        binding.NewFloat(),
		saturation: binding.NewFloat(),
		value:      binding.NewFloat(),
		unused:     binding.NewBool(),
	}

	pe.setFields()
	pe.patch.Editing = true

	hueLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(pe.hue, "%3.0f"))
	hueSlider := NewButtonSlide(pe.hue, HueBounds)
	hueBox := container.NewBorder(nil, nil, widget.NewLabel("H"), hueLabel,
		container.NewBorder(nil, nil, nil, nil, hueSlider.Container))

	saturationLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(pe.saturation, "%3.0f"))
	saturationSlider := NewButtonSlide(pe.saturation, SaturationBounds)
	saturationBox := container.NewBorder(nil, nil, widget.NewLabel("S"), saturationLabel,
		saturationSlider.Container)

	valueLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(pe.value, "%3.0f"))
	valueSlider := NewButtonSlide(pe.value, ValueBounds)
	valueBox := container.NewBorder(nil, nil, widget.NewLabel("V"), valueLabel,
		valueSlider.Container)

	pickerButton := widget.NewButtonWithIcon("Pick", theme.MoreHorizontalIcon(),
		pe.selectColorPicker(pe.patch))
	pe.removeButton = widget.NewButtonWithIcon("Cut", theme.ContentCutIcon(),
		pe.remove)

	pe.hue.AddListener(binding.NewDataListener(pe.setHue))
	pe.saturation.AddListener(binding.NewDataListener(pe.setSaturation))
	pe.value.AddListener(binding.NewDataListener(pe.setValue))

	revertButton := widget.NewButtonWithIcon("Cancel",
		theme.CancelIcon(), func() {
			pe.Dialog.Hide()
		})
	pe.applyButton = widget.NewButtonWithIcon("Apply",
		theme.ConfirmIcon(), pe.apply)
	vbox := container.NewVBox(
		pe.patch,
		hueBox,
		saturationBox,
		valueBox,
		widget.NewSeparator(), pickerButton)

	pe.Dialog = dialog.NewCustomWithoutButtons("", vbox, window)
	pe.Dialog.SetButtons([]fyne.CanvasObject{revertButton, pe.applyButton, pe.removeButton})
	return pe
}

func (pe *ColorPatchEditor) setFields() {
	pe.hsv.FromRGB(pe.patch.FillColor)
	fmt.Println("setfields", pe.hsv.Hue, pe.hsv.Saturation)
	pe.hue.Set(float64(pe.hsv.Hue))
	pe.saturation.Set(float64(pe.hsv.Saturation) * 100)
	pe.value.Set(float64(pe.hsv.Value) * 100)
	pe.unused.Set(pe.patch.unused)
}

func (pe *ColorPatchEditor) GetHSVColor() (hsv HSV) {
	hsv.FromRGB(pe.patch.FillColor)
	return
}

func (pe *ColorPatchEditor) setHue() {
	h, _ := pe.hue.Get()
	hsv := pe.GetHSVColor()
	hsv.Hue = float32(h)
	pe.setHSVColor(hsv)
}
func (pe *ColorPatchEditor) setSaturation() {
	s, _ := pe.saturation.Get()
	hsv := pe.GetHSVColor()
	hsv.Saturation = float32(s / 100)
	pe.setHSVColor(hsv)
}
func (pe *ColorPatchEditor) setValue() {
	v, _ := pe.value.Get()
	hsv := pe.GetHSVColor()
	hsv.Value = float32(v / 100)
	pe.setHSVColor(hsv)
}

func (pe *ColorPatchEditor) setHSVColor(hsv HSV) {
	pe.patch.FillColor = hsv.ToRGB()
}

func (pe *ColorPatchEditor) remove() {
	pe.patch.EditCut()
	pe.apply()
}

func (pe *ColorPatchEditor) apply() {
	pe.source.CopyPatch(pe.patch)
	pe.onUpdate()
	pe.Dialog.Hide()
}

func (le *ColorPatchEditor) selectColorPicker(patch *ColorPatch) func() {
	return func() {
		picker := dialog.NewColorPicker("Color Picker", "color", func(c color.Color) {
			if c != patch.GetColor() {
				r, g, b, a := c.RGBA()
				patch.FillColor = color.NRGBA{R: uint8(r), G: uint8(g),
					B: uint8(b), A: uint8(a)}
			}
		}, le.window)
		picker.Advanced = true
		picker.SetColor(patch.GetColor())
		picker.Show()
	}
}
