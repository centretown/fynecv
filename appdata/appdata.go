package appdata

import (
	"fynecv/entity"
	"fynecv/vision"

	"gocv.io/x/gocv"
)

const (
	V4LCAM = iota
	ESPCAM
)

type AppData struct {
	Cameras []*vision.Camera
	Lights  []*entity.Light
	Actions []*entity.Number
}

func NewAppData() *AppData {
	var data = &AppData{
		Cameras: []*vision.Camera{
			vision.NewCamera(0, gocv.VideoCaptureV4L),
			vision.NewCamera("http://192.168.0.28:8080", gocv.VideoCaptureAny),
		},
		Lights: entity.NewLightList([]string{
			"light.led_matrix_24",
			"light.led_strip_24"}),
		Actions: entity.NewNumberList([]string{
			"number.pan",
			"number.tilt"}),
	}
	return data
}
