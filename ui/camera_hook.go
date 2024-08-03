package ui

import (
	"fynecv/vision"
	"image"

	"fyne.io/fyne/v2/canvas"
)

type CameraHook struct {
	imageBox *canvas.Image
}

var _ vision.Hook = (*CameraHook)(nil)
var _ vision.UiHook = (*CameraHook)(nil)

func NewCameraHook(imageBox *canvas.Image) *CameraHook {
	fh := &CameraHook{
		imageBox: imageBox,
	}
	return fh
}

func (fh *CameraHook) CloseHook(int) {}

func (fh *CameraHook) UpdateHook(img image.Image) {
	if fh.imageBox != nil {
		fh.imageBox.Image = img
		fh.imageBox.Refresh()
	}
}

func (fh *CameraHook) SetUi(imageBox interface{}) {
	fh.imageBox = imageBox.(*canvas.Image)
}
