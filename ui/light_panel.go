package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/svc"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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

	patch := NewColorPatch()
	sel := widget.NewSelect([]string{}, nil)

	brightPct := binding.NewFloat()
	brightText := binding.FloatToStringWithFormat(brightPct, "%3.0f%%")
	brightLabel := widget.NewLabelWithData(brightText)
	slider := widget.NewSlider(0, 100)

	sel.OnChanged = func(s string) {
		if !fromSubscription {
			effect = s
			dataMap := make(map[string]string)
			dataMap["effect"] = fmt.Sprintf(`"%s"`, effect)
			cmd := svc.LightCmd(entityID, dataMap)
			_, err := data.CallService(cmd)
			if err != nil {
				log.Println(err, "CallService")
			}
		}
		fromSubscription = false
	}

	slider.OnChangeEnded = func(f float64) {
		if !fromSubscription {
			err := brightPct.Set(f)
			if err != nil {
				log.Println(err, "CallService")
			}
			v, err := brightPct.Get()
			if err != nil {
				log.Println(err, v, "CallService")
			}
			b, err := brightText.Get()
			if err != nil {
				log.Println(err, b, "CallService")
			}
			log.Println("slider", f, "pct", v, "pcttext", b, "CallService")

			brightValue = int(f) * 255 / 100
			dataMap := make(map[string]string)
			dataMap["brightness"] = fmt.Sprintf("%d", brightValue)
			cmd := svc.LightCmd(entityID, dataMap)
			_, err = data.CallService(cmd)
			if err != nil {
				log.Println(err, "CallService")
			}
		}
		fromSubscription = false
	}

	patch.SetOnTapped(func() {
		picker := dialog.NewColorPicker("Color Picker", "color", func(c color.Color) {
			rgb := nrgba(c)
			if red != rgb.R || green != rgb.G || blue != rgb.B {
				red, green, blue = rgb.R, rgb.G, rgb.B
				patch.SetColor(rgb)
				fmt.Println("patch.FillColor", patch.FillColor)
				dataMap := make(map[string]string)
				dataMap["rgb_color"] = fmt.Sprintf("[%d,%d,%d]", red, green, blue)
				cmd := svc.LightCmd(entityID, dataMap)
				_, err := data.CallService(cmd)
				if err != nil {
					log.Println(err, "CallService")
				}
			}
		}, win)
		picker.Advanced = true
		picker.SetColor(patch.GetColor())
		picker.Show()
	})

	data.Subscribe(entityID,
		appdata.NewSubcription(&light.Entity, func(c appdata.Consumer) {
			if light.Attributes.Brightness != brightValue {
				fromSubscription = true
				brightValue = light.Attributes.Brightness
				v := float64(brightValue) * 100 / 255

				brightPct.Set(v)
				slider.SetValue(v)
			}
			if light.Attributes.Effect != effect {
				fromSubscription = true
				effect = light.Attributes.Effect
				sel.Options = light.Attributes.EffectList
				sel.SetSelected(effect)
			}

			rgb := light.Attributes.ColorRGB
			if len(rgb) > 2 {
				if red != rgb[0] || green != rgb[1] || blue != rgb[2] {
					red, green, blue = rgb[0], rgb[1], rgb[2]
					patch.SetColor(color.NRGBA{R: red, G: green, B: blue, A: 255})
					patch.Refresh()
				}
			}
			log.Println("Refresh")
		}))

	return container.NewBorder(nil, nil,
		container.NewHBox(
			// widget.NewIcon(EffectIcon),
			widget.NewLabel("Effect:"),
			sel,
			widget.NewIcon(theme.ColorChromaticIcon()),
			widget.NewLabel("Color:"),
			patch,
			// widget.NewIcon(BrightIcon),
			brightLabel,
		),

		nil, slider)
}
