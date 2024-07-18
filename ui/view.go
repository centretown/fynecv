package ui

import (
	"fynecv/appdata"
	"fynecv/vision"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CameraView struct {
	Data       *appdata.AppData
	Container  *fyne.Container
	Image      *canvas.Image
	Current    int
	MainHook   vision.UiHook
	panSlider  *widget.Slider
	tiltSlider *widget.Slider
}

const (
	msgStartRecording = "Start Recording"
	msgStopRecording  = "Stop Recording"
)

func NewCameraView(data *appdata.AppData) *CameraView {

	cv := &CameraView{
		Data:  data,
		Image: canvas.NewImageFromImage(image.NewNRGBA(image.Rect(0, 0, 1280, 720))),
	}
	cv.Image.FillMode = canvas.ImageFillContain
	cv.MainHook = NewCameraHook(cv.Image)
	for _, camera := range data.Cameras {
		camera.MainHook = cv.MainHook
	}

	cv.panSlider = NewNumberSlider("number.pan", data)
	cv.tiltSlider = NewNumberSlider("number.tilt", data)
	cv.tiltSlider.Orientation = widget.Vertical

	if len(data.Cameras) > 0 {
		cam := data.Cameras[0]
		cam.HideMain = false
		cv.SetPanTilt(cam)
	}

	bottom := container.NewBorder(nil, nil, NewRecordButton(data), nil, cv.panSlider)

	cv.Container = container.NewBorder(
		nil,
		bottom,
		nil,
		cv.tiltSlider,
		cv.Image)

	return cv
}

func (cv *CameraView) SetCamera(id int) {
	cameras := cv.Data.Cameras
	if id != cv.Current {
		current := cameras[cv.Current]
		next := cameras[id]
		if next.Active {
			current.DisableMain()
			next.EnableMain()
			cv.Current = id
			cv.SetPanTilt(next)
		}
	}
}

func (cv *CameraView) SetPanTilt(cam *vision.Camera) {
	if cam.PanID == "" {
		cv.panSlider.Hide()
	} else {
		cv.panSlider.Show()
	}
	if cam.TiltID == "" {
		cv.tiltSlider.Hide()
	} else {
		cv.tiltSlider.Show()
	}
}
