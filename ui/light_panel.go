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

func NewLightPanel(entityID string, win fyne.Window, data *appdata.AppData) *fyne.Container {
	var (
		light            appdata.Light
		fromSubscription bool
		brightValue      int
		effect           string
		red, green, blue uint8
	)
	nrgba := func(c color.Color) color.NRGBA {
		r, g, b, a := c.RGBA()
		return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	}

	colorPatch := NewColorPatch()
	effectSelect := widget.NewSelect([]string{}, nil)
	brightSlider := widget.NewSlider(0, 100)
	brightLabel := widget.NewLabel("  0%")

	effectSelect.OnChanged = func(s string) {
		if !fromSubscription {
			effect = s
			data.CallService(svc.LightCmd(entityID, svc.ServiceData{
				Key:   "effect",
				Value: fmt.Sprintf(`"%s"`, effect),
			}))
		}
		fromSubscription = false
	}

	brightSlider.OnChangeEnded = func(f float64) {
		if !fromSubscription {
			brightLabel.SetText(fmt.Sprintf("%3.0f%%", f))
			brightValue = int(f) * 255 / 100
			data.CallService(svc.LightCmd(entityID, svc.ServiceData{
				Key:   "brightness",
				Value: fmt.Sprintf("%d", brightValue),
			}))
		}
		fromSubscription = false
	}

	colorPatch.SetOnTapped(func() {
		picker := dialog.NewColorPicker("Color Picker", "color", func(c color.Color) {
			rgb := nrgba(c)
			if red != rgb.R || green != rgb.G || blue != rgb.B {
				red, green, blue = rgb.R, rgb.G, rgb.B
				colorPatch.SetColor(rgb)
				data.CallService(svc.LightCmd(entityID, svc.ServiceData{
					Key:   "rgb_color",
					Value: fmt.Sprintf("[%d,%d,%d]", red, green, blue),
				}))
			}
		}, win)
		picker.Advanced = true
		picker.SetColor(colorPatch.GetColor())
		picker.Show()
	})

	data.Subscribe(entityID,
		appdata.NewSubcription(&light.Entity, func(c appdata.Consumer) {
			if light.Attributes.Brightness != brightValue {
				fromSubscription = true
				brightValue = light.Attributes.Brightness
				v := float64(brightValue) * 100 / 255

				brightLabel.SetText(fmt.Sprintf("%3.0f%%", v))
				brightSlider.SetValue(v)
			}
			if light.Attributes.Effect != effect {
				fromSubscription = true
				effect = light.Attributes.Effect
				effectSelect.Options = light.Attributes.EffectList
				effectSelect.SetSelected(effect)
			}

			rgb := light.Attributes.ColorRGB
			if len(rgb) > 2 {
				if red != rgb[0] || green != rgb[1] || blue != rgb[2] {
					red, green, blue = rgb[0], rgb[1], rgb[2]
					colorPatch.SetColor(color.NRGBA{R: red, G: green, B: blue, A: 255})
					colorPatch.Refresh()
				}
			}
			// log.Println("Refresh")
		}))

	return container.NewBorder(nil, nil,
		container.NewHBox(
			widget.NewIcon(EffectIcon),
			widget.NewLabel("Effect:"),
			effectSelect,
			widget.NewIcon(theme.ColorChromaticIcon()),
			widget.NewLabel("Color:"),
			colorPatch,
			widget.NewIcon(BrightIcon),
			brightLabel,
		),

		nil, brightSlider)
}
