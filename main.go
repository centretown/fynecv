package main

import (
	"flag"
	"log"
	"os"

	"fynecv/appdata"
	"fynecv/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	run(app, win)
}

func run(app fyne.App, win fyne.Window) {
	data := appdata.NewAppData()
	data.Monitor()

	wait := make(chan int)

	data.Ready.AddListener(binding.NewDataListener(func() {
		ready, _ := data.Ready.Get()
		if ready {
			log.Println("STATE LOADED")
		}

		if data.Err != nil {
			log.Fatal(data.Err)
			os.Exit(1)
		}

		wait <- 1
	}))

	<-wait
	view := ui.NewView(data)
	cameraList := ui.NewCameraList(data, win, view)
	cameraList.List.OnSelected = func(id widget.ListItemID) {
		view.SetCamera(id)
	}
	lightPanel := ui.NewPanel(data, win)
	ctr := container.NewBorder(lightPanel.Tabs,
		nil, nil, cameraList.Container, view.Container)
	win.SetContent(ctr)
	win.Resize(fyne.NewSize(1280+250, 768+100))
	win.Show()

	for _, camera := range data.Cameras {
		go camera.Serve()
	}

	win.SetCloseIntercept(func() {
		for _, device := range data.Cameras {
			if device.Busy {
				device.Quit <- 1
			}
		}
		win.Close()
	})

	app.Run()
}
