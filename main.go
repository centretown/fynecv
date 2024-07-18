package main

import (
	"flag"

	"fynecv/appdata"
	"fynecv/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const APPID = "com.centretown.fynecv.preferences"

func main() {
	flag.Parse()

	app := app.NewWithID(APPID)
	preferences := app.Preferences()
	theme := NewGlowTheme(preferences)
	app.Settings().SetTheme(theme)

	win := app.NewWindow("Cameras+Lights+Actions")

	data := setup(win)

	win.SetCloseIntercept(func() {
		for _, device := range data.Cameras {
			if device.Active {
				device.Quit <- 1
			}
		}
		win.Close()
	})

	win.Resize(fyne.NewSize(1280+250, 768+100))
	win.ShowAndRun()
}

func setup(win fyne.Window) *appdata.AppData {
	// data := appdata.NewAppData()
	data := appdata.NewAppData()

	view := ui.NewCameraView(data)
	tabs := ui.NewTabs(data, win)

	cameraList := ui.NewCameraList(data, win, view)
	cameraList.List.OnSelected = func(id widget.ListItemID) {
		view.SetCamera(id)
	}

	data.GetReady()

	// spl := container.NewHSplit(view.Container, cameraList.Container)

	ctr := container.NewBorder(tabs,
		nil, nil, cameraList.List, view.Container)
	win.SetContent(ctr)

	for _, camera := range data.Cameras {
		go camera.Serve()
	}
	return data
}
