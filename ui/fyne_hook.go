package ui

import (
	"fynecv/vision"

	"fyne.io/fyne/v2/canvas"
	"gocv.io/x/gocv"
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

func (fh *FyneHook) Update(mat *gocv.Mat) {
	if img, err := mat.ToImage(); err == nil {
		fh.imageBox.Image = img
		fh.imageBox.Refresh()
	}
}

func (fh *FyneHook) SetUi(imageBox interface{}) {
	fh.imageBox = imageBox.(*canvas.Image)
}
