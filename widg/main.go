package main

import (
	"flag"
	"fynecv/ui"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func main() {
	flag.Parse()

	app := app.New()
	win := app.NewWindow("WIDG")
	rect := image.Rect(0, 0, 1280, 720)
	img := image.NewNRGBA(rect)
	canv := canvas.NewImageFromImage(img)
	ai := ui.NewCameraWidget(canv)
	// ctr := container.NewBorder(widget.NewLabel("Hello World"), nil, nil, nil, canv)
	win.SetContent(ai)
	win.Resize(fyne.NewSize(1280, 800))
	win.ShowAndRun()
}
