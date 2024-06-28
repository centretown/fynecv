package cv

import (
	"log"
	"time"

	"gocv.io/x/gocv"
)

type Verb uint16

const (
	Get Verb = iota
	Set
	AddHook
	RemoveHook
)

var cmdList = []string{"GET", "SET"}

func (cmd Verb) String() string {
	if cmd > Set {
		return "Unknown"
	}
	return cmdList[cmd]
}

type CameraCmd struct {
	Action   Verb
	Property gocv.VideoCaptureProperties
	Value    any
}

type Camera struct {
	ID  any
	API gocv.VideoCaptureAPI

	Quit chan int
	Cmd  chan CameraCmd

	StreamHook *StreamHook
	ThumbHook  UiHook
	MainHook   UiHook
	Filters    []Hook

	capture     *gocv.VideoCapture
	framewidth  float64
	frameheight float64
	fps         float64
	ShowMain    bool
}

func NewCamera(id any, api gocv.VideoCaptureAPI) *Camera {
	dev := &Camera{
		ID:         id,
		API:        api,
		Quit:       make(chan int),
		Cmd:        make(chan CameraCmd),
		StreamHook: NewStreamHook(),
		Filters:    make([]Hook, 0),

		framewidth:  1280,
		frameheight: 720,
		fps:         20,
	}
	return dev
}

func (dev *Camera) AddFilter(filter Hook) {
	dev.Filters = append(dev.Filters, filter)
}

func (dev *Camera) RemoveMain() {
	devCmd := CameraCmd{Action: RemoveHook, Value: 0}
	dev.Cmd <- devCmd
}

func (dev *Camera) AddMain(hook Hook) {
	devCmd := CameraCmd{Action: AddHook, Value: hook}
	dev.Cmd <- devCmd
}

func (dev *Camera) do(cmd CameraCmd) {
	switch cmd.Action {
	case Get:
		cmd.Value = dev.capture.Get(cmd.Property)
	case Set:
		f, _ := cmd.Value.(float64)
		dev.capture.Set(cmd.Property, float64(f))
	case AddHook:
		dev.ShowMain = true
	case RemoveHook:
		dev.ShowMain = false
	}

}

func (dev *Camera) Open() (err error) {
	var (
		useAPI = dev.API > 0
	)
	if useAPI {
		dev.capture, err = gocv.OpenVideoCaptureWithAPI(dev.ID, dev.API)
	} else {
		dev.capture, err = gocv.OpenVideoCapture(dev.ID)
	}

	if err != nil {
		log.Println(err, dev.ID, "OpenVideoCapture")
		return
	}

	if useAPI {
		dev.capture.Set(gocv.VideoCaptureFPS, dev.fps)
		dev.capture.Set(gocv.VideoCaptureFrameHeight, dev.frameheight)
		dev.capture.Set(gocv.VideoCaptureFrameWidth, dev.framewidth)
	}

	dev.framewidth = dev.capture.Get(gocv.VideoCaptureFrameWidth)
	dev.frameheight = dev.capture.Get(gocv.VideoCaptureFrameHeight)
	dev.fps = dev.capture.Get(gocv.VideoCaptureFPS)
	log.Printf("Size: %.0fx%.0f FPS: %.0f\n", dev.framewidth, dev.frameheight, dev.fps)
	return
}

func (dev *Camera) Close() {
	// for i, hook := range dev.Hooks {
	// 	hook.Close(0)
	// 	log.Println("Closed hook", i)
	// }
	log.Println("Closed image")
	dev.capture.Close()
	log.Println("Closed capture")
	log.Println("Closed device", dev.ID)
}

func (dev *Camera) Capture() {
	var (
		img   gocv.Mat
		cmd   CameraCmd
		retry int = 0
	)

	err := dev.Open()
	if err != nil {
		return
	}

	close := func() {
		dev.Close()
		img.Close()
	}

	defer close()

	img = gocv.NewMat()

	for {
		time.Sleep(time.Millisecond * 25)

		select {

		case <-dev.Quit:
			return

		case cmd = <-dev.Cmd:
			dev.do(cmd)
			continue

		default:

		}

		if !dev.capture.Read(&img) {
			if retry > 10 {
				log.Println("Device unavailable:", dev.ID, retry)
				return
			}
			log.Println("Device closed:", dev.ID, retry)
			time.Sleep(100 * time.Millisecond)
			retry++
			dev.Open()
			continue
		}

		retry = 0

		if img.Empty() {
			continue
		}

		for _, filter := range dev.Filters {
			filter.Update(&img)
		}

		if dev.ShowMain {
			dev.MainHook.Update(&img)
		}
		dev.ThumbHook.Update(&img)
		dev.StreamHook.Update(&img)
	}

}
