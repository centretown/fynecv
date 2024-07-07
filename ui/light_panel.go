package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/comm"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Panel struct {
	data *appdata.AppData
	Tabs *container.AppTabs
}

func NewPanel(data *appdata.AppData, win fyne.Window) *Panel {
	lp := &Panel{
		data: data,
		Tabs: container.NewAppTabs(),
	}

	for _, light := range data.Lights {
		tab := container.NewTabItem(light.Attributes.Name, newLightContainer(light, win))
		lp.Tabs.Append(tab)
	}
	return lp
}

func newLightContainer(light *appdata.Light, win fyne.Window) *fyne.Container {

	sel := widget.NewSelect(light.Attributes.EffectList, func(s string) {
		comm.Post("services/light/turn_on",
			fmt.Sprintf(`{"entity_id": "%s", "effect": "%s"}`, light.EntityID, s))
	})
	sel.SetSelected(light.Attributes.Effect)

	brightBound := binding.NewFloat()
	brightValue := binding.NewSprintf("%.0f", brightBound)
	brightBound.Set(float64(light.Attributes.Brightness))
	brightBound.AddListener(binding.NewDataListener(func() {
		value, _ := brightBound.Get()
		comm.Post("services/light/turn_on",
			fmt.Sprintf(`{"entity_id": "%s", "brightness_pct": %.0f}`,
				light.EntityID, value))
	}))
	slider := widget.NewSliderWithData(0, 100, brightBound)

	var hsv HSV
	if len(light.Attributes.ColorHS) > 1 {
		hsv.Hue = float32(light.Attributes.ColorHS[0])
		hsv.Saturation = float32(light.Attributes.ColorHS[1])
		hsv.Value = float32(light.Attributes.Brightness)
	}

	patch := NewColorPatchWithColor(hsv, nil, nil)
	patch.SetOnTapped(func() {
		ce := NewColorPatchEditor(patch, win, func() {
			comm.Post("services/light/turn_on",
				fmt.Sprintf(`{"entity_id": "%s", "hs_color": [%f, %f]}`,
					light.EntityID,
					patch.colorHSV.Hue, patch.colorHSV.Saturation*100))
			brightBound.Set(float64(patch.colorHSV.Value * 100))
		})
		ce.Dialog.Show()
	})

	return container.NewBorder(nil, nil,
		container.NewHBox(
			container.NewHBox(widget.NewIcon(EffectIcon), widget.NewLabel("Effect:"), sel),
			container.NewHBox(widget.NewIcon(theme.ColorChromaticIcon()),
				widget.NewLabel("Color:"), patch)),
		nil,
		container.NewBorder(
			nil, nil,
			container.NewHBox(
				widget.NewIcon(BrightIcon),
				widget.NewLabel("Brightness:"),
				widget.NewLabelWithData(brightValue)),
			nil,
			slider))

}
