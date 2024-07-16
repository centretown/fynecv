package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/svc"
	"fynecv/vision"
	"image"
	"io"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type View struct {
	Data           *appdata.AppData
	Container      *fyne.Container
	Image          *canvas.Image
	Current        int
	MainHook       vision.UiHook
	IsRecording    bool
	RecordDuration int
}

const (
	msgStartRecording = "Start Recording"
	msgStopRecording  = "Stop Recording"
)

func NewView(data *appdata.AppData) *View {

	mv := &View{
		Data:           data,
		Image:          canvas.NewImageFromImage(image.NewNRGBA(image.Rect(0, 0, 1280, 720))),
		RecordDuration: 300,
	}

	var (
		panValue, tiltValue   float64
		panNumber, tiltNumber appdata.Number
	)

	fromSubscription := false
	onChangeEnded := func(f float64, entityID string) {
		if !fromSubscription {
			data.CallService(svc.NumberCmd(entityID, svc.ServiceData{
				Key:   "value",
				Value: fmt.Sprintf("%.0f", f),
			}))
		}
		fromSubscription = false

	}

	panSlider := widget.NewSlider(0, 180)
	panSlider.OnChangeEnded = func(f float64) {
		panValue = f
		onChangeEnded(f, "number.pan")
	}

	tiltSlider := widget.NewSlider(0, 180)
	tiltSlider.Orientation = widget.Vertical
	tiltSlider.OnChangeEnded = func(f float64) {
		tiltValue = f
		onChangeEnded(f, "number.tilt")
	}

	data.Subscribe("number.pan",
		appdata.NewSubcription(&panNumber.Entity, func(c appdata.Consumer) {
			var f float64
			_, err := fmt.Sscanf(panNumber.State, "%f", &f)
			if err != nil {
				log.Println(err, "pan state")
				return
			}
			if f != panValue {
				panValue = f
				fromSubscription = true
				panSlider.SetValue(panValue)
			}
		}))

	data.Subscribe("number.tilt",
		appdata.NewSubcription(&tiltNumber.Entity, func(c appdata.Consumer) {
			var f float64
			_, err := fmt.Sscanf(panNumber.State, "%f", &f)
			if err != nil {
				log.Println(err, "tilt state")
				return
			}
			if f != tiltValue {
				tiltValue = f
				fromSubscription = true
				tiltSlider.SetValue(tiltValue)
			}
		}))

	mv.Image.FillMode = canvas.ImageFillContain
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
			req := fmt.Sprintf("http://192.168.0.7:9000/record?duration=%d", mv.RecordDuration)
			resp, err = http.Get(req)
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

	bottom := container.NewBorder(nil, nil, recordButton, nil, panSlider)

	mv.Container = container.NewBorder(
		nil,
		bottom,
		nil,
		tiltSlider,
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

// func (view *View) Bind() {
// 	view.Binder.Pan.Set(140)
// 	// view.Binder.Tilt.Set(50)

// 	view.Binder.Pan.AddListener(binding.NewDataListener(func() {
// 		value, _ := view.Binder.Pan.Get()
// 		cmd := fmt.Sprintf(`{"entity_id": "number.pan", "value": %.0f}`, value)
// 		comm.Post("services/number/set_value", cmd)
// 	}))
// 	view.Binder.Tilt.AddListener(binding.NewDataListener(func() {
// 		value, _ := view.Binder.Tilt.Get()
// 		cmd := fmt.Sprintf(`{"entity_id": "number.tilt", "value": %.0f}`, value)
// 		comm.Post("services/number/set_value", cmd)
// 	}))

// }

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
