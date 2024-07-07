package ui

import (
	"fynecv/vision"
	"image"

	"fyne.io/fyne/v2/canvas"
)

type FyneHook struct {
	imageBox *canvas.Image
}

var _ vision.Hook = (*FyneHook)(nil)
var _ vision.UiHook = (*FyneHook)(nil)

func NewFyneHook(imageBox *canvas.Image) *FyneHook {
	fh := &FyneHook{
		imageBox: imageBox,
	}
	return fh
}

func (fh *FyneHook) Close(int) {}

func (fh *FyneHook) Update(img image.Image) {
	if fh.imageBox != nil {
		fh.imageBox.Image = img
		fh.imageBox.Refresh()
	}
}

func (fh *FyneHook) SetUi(imageBox interface{}) {
	fh.imageBox = imageBox.(*canvas.Image)
}
