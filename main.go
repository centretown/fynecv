package main

import (
	"flag"
	"fmt"
	"image"
	"log"

	"fynecv/cv"
	"fynecv/hass"
	"fynecv/web"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"gocv.io/x/gocv"
)

const APPID = "com.centretown.fynecv.preferences"

func main() {
	flag.Parse()
	app := app.NewWithID(APPID)
	preferences := app.Preferences()
	theme := NewGlowTheme(preferences)
	app.Settings().SetTheme(theme)
	win := app.NewWindow("Cameras")
	run(win)
}

func run(win fyne.Window) {
	var devices = []*cv.Device{
		cv.NewDevice(0, gocv.VideoCaptureV4L),
		cv.NewDevice("http://192.168.0.25:8080", gocv.VideoCaptureAny),
	}

	mainImage := canvas.NewImageFromImage(image.NewNRGBA(image.Rect(0, 0, 1280, 720)))
	mainImage.FillMode = canvas.ImageFillContain
	mainHook := NewFyneHook(mainImage)
	list := newList(devices, mainHook)

	ctr := container.NewBorder(nil, nil, nil, list, newView(mainImage))
	win.SetContent(ctr)

	go web.Serve(devices)

	win.SetCloseIntercept(func() {
		for _, device := range devices {
			device.Quit <- 1
		}
		win.Close()
	})

	win.Resize(fyne.NewSize(1280, 960))
	win.ShowAndRun()
}

func newList(devices []*cv.Device, mainHook cv.UiHook) fyne.CanvasObject {

	bound := binding.NewUntypedList()
	// classify := cv.NewClassifyHook()

	for _, device := range devices {
		device.MainHook = mainHook
		device.ThumbHook = NewFyneHook(nil)
		device.StreamHook = cv.NewStreamHook()
		// device.AddFilter(classify)
		bound.Append(device)
	}

	list := widget.NewListWithData(
		bound,
		func() fyne.CanvasObject {
			imageBox := canvas.NewImageFromImage(image.NewNRGBA(image.Rect(0, 0, 240, 200)))
			imageBox.FillMode = canvas.ImageFillContain
			imageBox.SetMinSize(fyne.NewSize(240, 200))
			return imageBox
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			d, _ := i.(binding.Untyped).Get()
			device, _ := d.(*cv.Device)
			device.ThumbHook.SetUi(o)
		},
	)

	current := 0
	devices[current].ShowMain = true

	list.OnSelected = func(id widget.ListItemID) {
		if id != current {
			log.Println("selected", id)
			devices[current].RemoveMain()
			current = id
			devices[current].AddMain(mainHook)
		}
	}

	return list
}

type ViewData struct {
	panValue  binding.Float
	tiltValue binding.Float
}

var viewData ViewData

func newView(view fyne.CanvasObject) fyne.CanvasObject {
	viewData.panValue = binding.NewFloat()
	viewData.panValue.Set(90)
	viewData.tiltValue = binding.NewFloat()
	viewData.tiltValue.Set(90)
	pan := widget.NewSliderWithData(0, 180, viewData.panValue)
	tilt := widget.NewSliderWithData(0, 180, viewData.tiltValue)
	tilt.Orientation = widget.Vertical
	ctr := container.NewBorder(nil,
		pan,
		nil,
		tilt,
		view)

	viewData.panValue.AddListener(binding.NewDataListener(func() {
		value, _ := viewData.panValue.Get()
		cmd := fmt.Sprintf(`{"entity_id": "number.pan", "value": %.0f}`, value)
		hass.Post("services/number/set_value", cmd)
		log.Println(value)
	}))
	viewData.tiltValue.AddListener(binding.NewDataListener(func() {
		value, _ := viewData.tiltValue.Get()
		cmd := fmt.Sprintf(`{"entity_id": "light.led_matrix_24","effect": "rainbow-vertical","brightness_pct": %.0f}`, value)
		hass.Post("services/light/turn_on", cmd)
		log.Println(value)
	}))
	return ctr
}
