package ui

import (
	"fynecv/cv"
	"image"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewCameraList(cameras []*cv.Camera, mainHook cv.UiHook) fyne.CanvasObject {
	const (
		thumbWidth  = 320
		thumbHeight = 200
	)

	bound := binding.NewUntypedList()
	// classify := cv.NewClassifyHook()

	for _, device := range cameras {
		device.MainHook = mainHook
		device.ThumbHook = NewFyneHook(nil)
		device.StreamHook = cv.NewStreamHook()
		// device.AddFilter(classify)
		bound.Append(device)
	}

	list := widget.NewListWithData(
		bound,
		func() fyne.CanvasObject {
			imageBox := canvas.NewImageFromImage(image.NewNRGBA(
				image.Rect(0, 0, thumbWidth, thumbHeight)))
			imageBox.FillMode = canvas.ImageFillContain
			imageBox.SetMinSize(fyne.NewSize(thumbWidth, thumbHeight))
			return imageBox
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			d, _ := i.(binding.Untyped).Get()
			device, _ := d.(*cv.Camera)
			device.ThumbHook.SetUi(o)
		},
	)

	current := 0
	cameras[current].ShowMain = true

	list.OnSelected = func(id widget.ListItemID) {
		if id != current {
			log.Println("selected", id)
			cameras[current].RemoveMain()
			current = id
			cameras[current].AddMain(mainHook)
		}
	}

	return list
}
