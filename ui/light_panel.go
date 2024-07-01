package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/hass"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type LightPanel struct {
	data *appdata.AppData
	Tabs *container.AppTabs
}

func NewLightPanel(data *appdata.AppData) *LightPanel {
	lp := &LightPanel{
		data: data,
		Tabs: container.NewAppTabs(),
	}

	for _, light := range data.Lights {

		sel := widget.NewSelect(light.Attributes.EffectList, func(s string) {
			cmd := fmt.Sprintf(`{"entity_id": "%s", "effect": "%s"}`,
				light.EntityID, s)
			hass.Post("services/light/turn_on", cmd)
		})
		sel.SetSelected(light.Attributes.Effect)

		brightness := binding.NewFloat()
		brightnessLabel := binding.NewSprintf("%.0f", brightness)
		brightness.Set(light.Attributes.Brightness)

		slider := widget.NewSliderWithData(0, 100, brightness)
		brightness.AddListener(binding.NewDataListener(func() {
			value, _ := brightness.Get()
			cmd := fmt.Sprintf(`{"entity_id": "%s", "brightness_pct": %.0f}`,
				light.EntityID, value)
			hass.Post("services/light/turn_on", cmd)
		}))

		ctr := container.NewBorder(nil, nil,
			container.NewHBox(widget.NewLabel("Effect"), sel), nil,
			container.NewBorder(nil, nil,
				container.NewHBox(widget.NewLabel("Brightness"),
					widget.NewLabelWithData(brightnessLabel)),
				nil, slider))
		tab := container.NewTabItem(light.Attributes.Name, ctr)

		lp.Tabs.Append(tab)
	}
	return lp
}

// 	lights := lp.data.Lights
// 	options = make([]string, 0, len(lights))

// 	for _, light := range lights {
// 		options = append(options, light.Attributes.Name)
// 	}

// 	return
// }

// func (lp *LightPanel) EffectList() (options []string) {
// 	effects := lp.data.Lights[lp.Current].Attributes.EffectList
// 	options = make([]string, 0, len(effects))
// 	options = append(options, effects...)
// 	return
// }

// func (lp *LightPanel) onChangeEffect(s string) {
// 	lp.Sets[lp.Current].Effect = s
// 	cmd := fmt.Sprintf(`{"entity_id": "%s", "effect": "%s"}`,
// 		lp.data.Lights[lp.Current].EntityID, lp.Effects.Selected)
// 	hass.Post("services/light/turn_on", cmd)
// }

// func (lp *LightPanel) onChange(s string) {
// 	lp.Current = lp.Select.SelectedIndex()
// 	set := lp.Sets[lp.Current]
// 	lp.Brightness.Set(set.Brightness)
// 	lp.Effects.Options = set.Effects
// 	lp.Effects.SetSelected(set.Effect)
// }

// func (lp *LightPanel) onSlide(s string) {
// 	value, _ := lp.Brightness.Get()
// 	lp.Sets[lp.Current].Brightness = value
// 	cmd := fmt.Sprintf(`{"entity_id": "%s", "brightness_pct": %.0f}`,
// 		lp.data.Lights[lp.Current].EntityID, value)
// 	hass.Post("services/light/turn_on", cmd)
// }
