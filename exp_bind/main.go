package main

import (
	"flag"
	"fmt"
	"fynecv/appdata"
	"fynecv/svc"
	"fynecv/ui"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// const (
// 	entityID = "light.led_matrix_24"
// )

func main() {
	flag.Parse()

	var (
		myapp    = app.New()
		win      = myapp.NewWindow("Conversion")
		data     = appdata.NewAppData()
		lightCtr = NewLightContainer(data, "light.led_strip_24")
	)

	button := widget.NewButton("Click", func() {
		log.Println("click")
		w2 := myapp.NewWindow("Lights")
		w2.SetContent(lightCtr)
		w2.Resize(fyne.NewSize(720, 150))
		w2.Show()
	})
	win.SetContent(button)
	data.Monitor()
	win.Resize(fyne.NewSize(720, 150))
	win.ShowAndRun()
}

type LightItem struct {
	Light      appdata.Light
	Subscribed bool
}

var idMap = map[string]*LightItem{
	"light.led_strip_24":  {},
	"light.led_matrix_24": {},
}

func NewTabs(data *appdata.AppData) *container.AppTabs {
	tabs := container.NewAppTabs()

	for entityID, item := range idMap {
		data.Subscribe(entityID,
			appdata.NewSubcription(&item.Light.Entity,
				func(c appdata.Consumer) {
					item, ok := idMap[entityID]
					if ok && item.Subscribed {
						return
					}
					idMap[entityID].Subscribed = true
					tabItem := container.NewTabItem(item.Light.Attributes.Name,
						NewLightContainer(data, entityID))
					tabs.Append(tabItem)
					tabs.Refresh()
				}))
	}

	return tabs
}

func NewLightContainer(data *appdata.AppData, entityID string) *fyne.Container {
	var (
		light   appdata.Light
		fromSub bool
		effect  string
	)

	brightValue := 0
	brightPct := binding.NewFloat()
	short := binding.FloatToStringWithFormat(brightPct, "%3.0f%%")
	slider := widget.NewSlider(0, 100)
	sel := widget.NewSelect([]string{}, nil)

	fill := color.NRGBA{A: 255}
	patch := ui.NewColorPatch()
	patch.FillColor = fill

	slider.OnChangeEnded = func(f float64) {
		if !fromSub {
			brightPct.Set(f)
			brightValue = int(f) * 255 / 100
			dataMap := make(map[string]string)
			dataMap["brightness"] = fmt.Sprintf("%d", brightValue)
			cmd := svc.LightCmd(entityID, dataMap)
			_, err := data.CallService(cmd)
			if err != nil {
				log.Println(err, "CallService")
			}
		}
		fromSub = false
	}

	sel.OnChanged = func(s string) {
		if !fromSub {
			effect = s
			dataMap := make(map[string]string)
			dataMap["effect"] = fmt.Sprintf(`"%s"`, effect)
			cmd := svc.LightCmd(entityID, dataMap)
			_, err := data.CallService(cmd)
			if err != nil {
				log.Println(err, "CallService")
			}
		}
		fromSub = false
	}

	ctr := container.NewBorder(
		nil, nil,
		container.NewHBox(
			widget.NewIcon(ui.EffectIcon), widget.NewLabel("Effect:"), sel,
			widget.NewIcon(theme.ColorChromaticIcon()),
			widget.NewLabel("Color:"),
			patch,
			widget.NewIcon(ui.BrightIcon),
		), nil,
		container.NewBorder(nil, nil,
			widget.NewLabelWithData(short), nil, slider))

	data.Subscribe(entityID,
		appdata.NewSubcription(&light.Entity, func(c appdata.Consumer) {
			if light.Attributes.Brightness != brightValue {
				fromSub = true
				brightValue = light.Attributes.Brightness
				v := float64(brightValue) * 100 / 255
				brightPct.Set(v)
				slider.SetValue(v)
			}
			if light.Attributes.Effect != effect {
				fromSub = true
				effect = light.Attributes.Effect
				sel.Options = light.Attributes.EffectList
				sel.SetSelected(effect)
			}

			rgb := light.Attributes.ColorRGB
			if len(rgb) > 2 {
				patch.SetColor(color.NRGBA{R: rgb[0], G: rgb[1], B: rgb[2], A: 255})
			}
			patch.Refresh()
			log.Println("Refresh")
		}))
	return ctr
}
