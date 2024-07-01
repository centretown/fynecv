package main

import (
	"flag"

	"fynecv/appdata"
	"fynecv/ui"
	"fynecv/web"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	data := appdata.NewAppData()
	view := NewView(data)
	mainHook := ui.NewFyneHook(view.Image)
	cameraList := ui.NewCameraList(data.Cameras, mainHook)
	lightPanel := ui.NewLightPanel(data)
	// accord := ui.NewAccordianList(lightList)

	// vbox := container.NewVBox(cameraList, lightList)

	ctr := container.NewBorder(lightPanel.Tabs,
		nil, nil, cameraList, view.Container)
	win.SetContent(ctr)

	go web.Serve(data.Cameras)

	win.SetCloseIntercept(func() {
		for _, device := range data.Cameras {
			device.Quit <- 1
		}
		win.Close()
	})

	win.Resize(fyne.NewSize(1280+250, 768+100))
	win.ShowAndRun()
}
