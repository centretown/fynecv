package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/comm"
	"image/color"

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
	appdata.ShowYaml(light)

	sel := widget.NewSelect(light.Attributes.EffectList, func(s string) {
		comm.Post("services/light/turn_on",
			fmt.Sprintf(`{"entity_id": "%s", "effect": "%s"}`, light.EntityID, s))
	})
	sel.SetSelected(light.Attributes.Effect)

	brightBound := binding.NewFloat()
	slider := widget.NewSliderWithData(0, 255, brightBound)
	brightBound.Set(float64(light.Attributes.Brightness))
	v, _ := brightBound.Get()
	fmt.Println("bound brightness", v)
	brightBound.AddListener(binding.NewDataListener(func() {
		v, _ := brightBound.Get()
		fmt.Println("bound brightness", v)
	}))

	brightValue := binding.NewSprintf("%.0f", brightBound)
	// brightBound.AddListener(binding.NewDataListener(func() {
	// 	value, _ := brightBound.Get()
	// 	fmt.Println("bound value", value)
	// 	comm.Post("services/light/turn_on",
	// 		fmt.Sprintf(`{"entity_id": "%s", "brightness_pct": %.0f}`,
	// 			light.EntityID, value))
	// }))

	var (
		hsv HSV
		rgb color.NRGBA
	)
	if len(light.Attributes.ColorRGB) > 2 {
		rgb.R = light.Attributes.ColorRGB[0]
		rgb.G = light.Attributes.ColorRGB[1]
		rgb.B = light.Attributes.ColorRGB[2]
		rgb.A = 255

		hsv.FromColor(rgb)
		appdata.ShowYaml(rgb)
		appdata.ShowYaml(hsv)

		// hsv.Hue = float32(light.Attributes.ColorHS[0])
		// hsv.Saturation = float32(light.Attributes.ColorHS[1]) / 100
		// hsv.Value = float32(light.Attributes.Brightness) / 255 * 100
	}

	patch := NewColorPatchWithColor(rgb, nil, nil)
	patch.SetOnTapped(func() {
		ce := NewColorPatchEditor(patch, win, func() {
			rgb := patch.ColorRGB
			comm.Post("services/light/turn_on",
				fmt.Sprintf(`{"entity_id": "%s", "rgb_color": [%d, %d, %d]}`,
					light.EntityID,
					rgb.R, rgb.G, rgb.B))
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
