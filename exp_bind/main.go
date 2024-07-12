package main

import (
	"flag"
	"fmt"
	"fynecv/appdata"
	"fynecv/ui"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var (
	serviceFormat = `{ "type":"call_service", "domain":"%s", "service":"%s", "service_data": {%s}, "target":{"entity_id":"%s"},`
	idstr         = `"id":%d }`
)

func serviceCmd(domain string, service string, entityID string, data map[string]string) string {
	var serviceData string
	for k, v := range data {
		serviceData += fmt.Sprintf(`"%s":%s,`, k, v)
	}
	if len(serviceData) > 0 {
		serviceData = serviceData[:len(serviceData)-1]
	}
	cmd := fmt.Sprintf(serviceFormat, domain, service, serviceData, entityID) + idstr
	return cmd
}

func main() {
	flag.Parse()
	data := appdata.NewAppData()

	myApp := app.New()
	w := myApp.NewWindow("Conversion")

	var (
		brightValue int
		// brightPct   float64
	)
	brightPct := binding.NewFloat()
	short := binding.FloatToStringWithFormat(brightPct, "%3.0f%%")
	slider := widget.NewSliderWithData(0, 100, brightPct)

	fill := color.NRGBA{A: 255}
	patch := ui.NewColorPatch()
	patch.ColorRGB = fill

	w.SetContent(container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabelWithData(short), nil, slider),
		patch,
	))

	var (
		light    appdata.Light
		fromSub  bool
		entityID = "light.led_matrix_24"
	)

	data.Subscribe(entityID,
		appdata.NewSubcription(&light.Entity, func(c appdata.Consumer) {
			if light.Attributes.Brightness != brightValue {
				fromSub = true
				brightValue = light.Attributes.Brightness
				slider.SetValue(float64(brightValue) * 100 / 255)
			}

			rgb := light.Attributes.ColorRGB
			if len(rgb) > 2 {
				patch.SetColor(color.NRGBA{R: rgb[0], G: rgb[1], B: rgb[2], A: 255})
			}
			patch.Refresh()
			log.Println("Refresh")
		}))

	slider.OnChangeEnded = func(f float64) {
		if !fromSub {
			// f, _ := brightPct.Get()
			brightValue = int(f) * 255 / 100
			d := make(map[string]string)
			d["brightness"] = fmt.Sprintf("%d", brightValue)
			cmd := serviceCmd("light", "turn_on", entityID, d)
			id, err := data.CallService(cmd)
			if err != nil {
				log.Println(err, "CallService")
			} else {
				log.Println(id, "CallService")
			}
		}
		fromSub = false
	}

	data.Monitor()
	w.Resize(fyne.NewSize(250, 250))
	w.ShowAndRun()
}
