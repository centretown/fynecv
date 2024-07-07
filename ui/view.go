package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/comm"
	"fynecv/vision"
	"image"
	"io"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type BindingData struct {
	Pan  binding.Float
	Tilt binding.Float
}

type View struct {
	Data        *appdata.AppData
	Binder      BindingData
	Container   *fyne.Container
	Image       *canvas.Image
	Pan         *widget.Slider
	Tilt        *widget.Slider
	Current     int
	MainHook    vision.UiHook
	IsRecording bool
}

const (
	msgStartRecording = "Start Recording"
	msgStopRecording  = "Stop Recording"
)

func NewView(data *appdata.AppData) *View {

	mv := &View{
		Data: data,
		Binder: BindingData{
			Pan:  binding.NewFloat(),
			Tilt: binding.NewFloat(),
		},
		Image: canvas.NewImageFromImage(image.NewNRGBA(image.Rect(0, 0, 1280, 720))),
	}

	mv.Pan = widget.NewSliderWithData(0, 180, mv.Binder.Pan)
	mv.Tilt = widget.NewSliderWithData(0, 180, mv.Binder.Tilt)
	mv.Bind()

	mv.Image.FillMode = canvas.ImageFillContain
	mv.Tilt.Orientation = widget.Vertical
	recordButton := widget.NewButtonWithIcon(msgStartRecording, MotionOnIcon, func() {})

	recordButton.OnTapped = func() {
		var (
			err  error
			resp *http.Response
		)

		if mv.IsRecording {
			recordButton.SetIcon(MotionOnIcon)
			recordButton.SetText(msgStartRecording)
			resp, err = http.Get("http://192.168.0.7:9000/record?duration=0")
		} else {
			recordButton.SetIcon(MotionOffIcon)
			recordButton.SetText(msgStopRecording)
			resp, err = http.Get("http://192.168.0.7:9000/record?duration=60")
		}

		if err != nil {
			log.Println("http.Get", err)
			return
		}

		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("ReadAll", err)
			return
		}

		log.Println(string(buf))
		mv.IsRecording = !mv.IsRecording
		recordButton.Refresh()
	}

	bottom := container.NewBorder(nil, nil, recordButton, nil, mv.Pan)

	mv.Container = container.NewBorder(
		nil,
		bottom,
		nil,
		mv.Tilt,
		mv.Image)

	mv.MainHook = NewFyneHook(mv.Image)

	for _, camera := range data.Cameras {
		camera.MainHook = mv.MainHook
	}

	if len(data.Cameras) > 0 {
		data.Cameras[0].HideMain = false
	}
	return mv
}

func (view *View) Bind() {
	view.Binder.Pan.Set(140)
	// view.Binder.Tilt.Set(50)

	view.Binder.Pan.AddListener(binding.NewDataListener(func() {
		value, _ := view.Binder.Pan.Get()
		cmd := fmt.Sprintf(`{"entity_id": "number.pan", "value": %.0f}`, value)
		comm.Post("services/number/set_value", cmd)
	}))
	view.Binder.Tilt.AddListener(binding.NewDataListener(func() {
		value, _ := view.Binder.Tilt.Get()
		cmd := fmt.Sprintf(`{"entity_id": "number.tilt", "value": %.0f}`, value)
		comm.Post("services/number/set_value", cmd)
	}))

}

func (view *View) SetCamera(id int) {
	cameras := view.Data.Cameras
	if id != view.Current {
		current := cameras[view.Current]
		next := cameras[id]
		if next.Busy {
			current.HideMainCmd()
			next.ShowMainCmd()
			view.Current = id
		}
	}
}
