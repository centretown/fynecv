package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/svc"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LightPanel struct {
	light appdata.Light

	effectLabel  *widget.Label
	effectIcon   *widget.Icon
	effectSelect *widget.Select

	colorLabel *widget.Label
	colorIcon  *widget.Icon
	colorPatch *ColorPatch

	brightLabel *widget.Label
	brightIcon  *widget.Icon
	brightValue *widget.Label

	brightSlider *widget.Slider
}

func NewLightPanel(entityID string, win fyne.Window, data *appdata.AppData) *fyne.Container {
	lp := &LightPanel{
		effectLabel:  widget.NewLabel("Effect:"),
		effectIcon:   widget.NewIcon(EffectIcon),
		effectSelect: widget.NewSelect([]string{}, nil),

		colorLabel: widget.NewLabel("Color:"),
		colorIcon:  widget.NewIcon(theme.ColorChromaticIcon()),
		colorPatch: NewColorPatch(),

		brightLabel:  widget.NewLabel("Brightness:"),
		brightIcon:   widget.NewIcon(BrightIcon),
		brightValue:  widget.NewLabel("  0%"),
		brightSlider: widget.NewSlider(0, 100),
	}

	var (
		fromSubscription bool
		brightValue      int
		effect           string
		red, green, blue uint8
	)
	nrgba := func(c color.Color) color.NRGBA {
		r, g, b, a := c.RGBA()
		return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	}

	lp.effectSelect.OnChanged = func(s string) {
		if !fromSubscription {
			effect = s
			data.CallService(svc.LightCmd(entityID, svc.ServiceData{
				Key:   "effect",
				Value: fmt.Sprintf(`"%s"`, effect),
			}))
		}
		fromSubscription = false
	}

	lp.brightSlider.OnChangeEnded = func(f float64) {
		if !fromSubscription {
			lp.brightValue.SetText(fmt.Sprintf("%3.0f%%", f))
			brightValue = int(f) * 255 / 100
			data.CallService(svc.LightCmd(entityID, svc.ServiceData{
				Key:   "brightness",
				Value: fmt.Sprintf("%d", brightValue),
			}))
		}
		fromSubscription = false
	}

	lp.colorPatch.SetOnTapped(func() {
		picker := dialog.NewColorPicker("Color Picker", "color", func(c color.Color) {
			rgb := nrgba(c)
			if red != rgb.R || green != rgb.G || blue != rgb.B {
				red, green, blue = rgb.R, rgb.G, rgb.B
				lp.colorPatch.SetColor(rgb)
				data.CallService(svc.LightCmd(entityID, svc.ServiceData{
					Key:   "rgb_color",
					Value: fmt.Sprintf("[%d,%d,%d]", red, green, blue),
				}))
			}
		}, win)
		picker.Advanced = true
		picker.SetColor(lp.colorPatch.GetColor())
		picker.Show()
	})

	data.Subscribe(entityID,
		appdata.NewSubcription(&lp.light.Entity, func(c appdata.Consumer) {
			if lp.light.Attributes.Brightness != brightValue {
				fromSubscription = true
				brightValue = lp.light.Attributes.Brightness
				v := float64(brightValue) * 100 / 255

				lp.brightValue.SetText(fmt.Sprintf("%3.0f%%", v))
				lp.brightSlider.SetValue(v)
			}
			if lp.light.Attributes.Effect != effect {
				fromSubscription = true
				effect = lp.light.Attributes.Effect
				lp.effectSelect.Options = lp.light.Attributes.EffectList
				lp.effectSelect.SetSelected(effect)
			}

			rgb := lp.light.Attributes.ColorRGB
			if len(rgb) > 2 {
				if red != rgb[0] || green != rgb[1] || blue != rgb[2] {
					red, green, blue = rgb[0], rgb[1], rgb[2]
					lp.colorPatch.SetColor(color.NRGBA{R: red, G: green, B: blue, A: 255})
					lp.colorPatch.Refresh()
				}
			}
		}))

	return lp.LayoutHorizontal()
}

func (lp *LightPanel) LayoutHorizontal() *fyne.Container {
	hbox := container.NewHBox(
		widget.NewSeparator(),
		container.NewHBox(
			lp.effectLabel,
			lp.effectIcon,
			lp.effectSelect),
		widget.NewSeparator(),
		container.NewHBox(
			lp.colorLabel,
			lp.colorIcon,
			lp.colorPatch),
		widget.NewSeparator(),
		container.NewHBox(
			lp.brightLabel,
			lp.brightIcon,
			lp.brightValue))
	ctr := container.NewBorder(nil, nil,
		hbox,
		nil, lp.brightSlider)

	return ctr

}
