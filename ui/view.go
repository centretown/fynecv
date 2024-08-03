package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/vision"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CameraView struct {
	Data         *appdata.AppData
	Container    *fyne.Container
	Image        *canvas.Image
	cameraWidget *CameraWidget
	Current      int
	panSlider    *widget.Slider
	tiltSlider   *widget.Slider
	MainHook     vision.UiHook
	viewStates   []*CameraWidgetwState
}

const (
	msgStartRecording = "Start Recording"
	msgStopRecording  = "Stop Recording"
)

func NewCameraView(data *appdata.AppData) *CameraView {

	cv := &CameraView{
		Current:    -1,
		Data:       data,
		Image:      canvas.NewImageFromImage(image.NewNRGBA(image.Rect(0, 0, 1280, 720))),
		viewStates: make([]*CameraWidgetwState, 0, len(data.Cameras)),
	}
	cv.Image.FillMode = canvas.ImageFillContain
	cv.cameraWidget = NewCameraWidget(cv.Image)
	cv.MainHook = cv.cameraWidget
	for _, camera := range data.Cameras {
		camera.MainHook = cv.cameraWidget
		cv.viewStates = append(cv.viewStates, &CameraWidgetwState{})
	}

	cv.panSlider = NewNumberSlider("number.pan", data)
	cv.tiltSlider = NewNumberSlider("number.tilt", data)
	cv.tiltSlider.Orientation = widget.Vertical

	if len(data.Cameras) > 0 {
		data.Cameras[0].HideMain = false
		cv.SetCamera(0)
	}

	bottom := container.NewBorder(nil, nil, NewRecordButton(data), nil, cv.panSlider)

	cv.Container = container.NewBorder(
		nil,
		bottom,
		nil,
		cv.tiltSlider,
		cv.cameraWidget)

	return cv
}

func (cv *CameraView) SetCamera(id int) {
	if id == cv.Current {
		return
	}

	cameras := cv.Data.Cameras
	if id >= len(cameras) || id < 0 {
		return
	}

	if cv.Current >= 0 {
		cameras[cv.Current].DisableMain()
	}

	camera := cameras[id]
	cv.Current = id
	cv.cameraWidget.SetState(cv.viewStates[id])

	if !camera.Active {
		fmt.Println("SetCamera", "!camera.Active", id)
		return
	}

	camera.EnableMain()
	cv.SetPanTilt(camera)

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
