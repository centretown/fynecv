package vision

import (
	"log"
	"time"

	"github.com/mattn/go-mjpeg"
)

type Verb uint16

const (
	GET Verb = iota
	SET
	HIDEMAIN
	HIDETHUMB
	HIDEALL
	RECORD_START
	RECORD_STOP
)

const (
	RecordingFolder = "recordings/"
)

var cmdList = []string{
	"Get",
	"Set",
	"HideMain",
	"HideThumb",
	"HideAll",
}

func (cmd Verb) String() string {
	if cmd >= Verb(len(cmdList)) {
		return "Unknown"
	}
	return cmdList[cmd]
}

type CameraCmd struct {
	Action Verb
	Value  any
}

type Camera struct {
	URL     string
	Name    string
	Decoder *mjpeg.Decoder

	HasPan  bool
	HasTilt bool

	Quit   chan int
	Record chan int
	Cmd    chan CameraCmd

	recordStop time.Time

	ThumbHook UiHook
	MainHook  UiHook

	Filters []Hook

	HideMain  bool
	HideThumb bool
	HideAll   bool
	Busy      bool
	Recording bool

	FrameWidth  float64
	FrameHeight float64
	FrameRate   float64
}

func NewCamera(url string) *Camera {

	cam := &Camera{
		URL:     url,
		Name:    url,
		Quit:    make(chan int),
		Record:  make(chan int),
		Cmd:     make(chan CameraCmd),
		Filters: make([]Hook, 0),

		FrameWidth:  1280,
		FrameHeight: 720,
		FrameRate:   20,
		HideMain:    true,
		HideThumb:   false,
		// writer:      &gocv.VideoWriter{},
	}

	return cam
}

func (cam *Camera) AddFilter(filter Hook) {
	cam.Filters = append(cam.Filters, filter)
}
func (cam *Camera) Command(cmd CameraCmd) {
	cam.Cmd <- cmd
}

func (cam *Camera) RecordCmd() {
	cam.Command(CameraCmd{Action: RECORD_START, Value: true})
}

func (cam *Camera) StopRecordCmd() {
	cam.Command(CameraCmd{Action: RECORD_STOP, Value: true})
}

func (cam *Camera) HideMainCmd() {
	cam.Command(CameraCmd{Action: HIDEMAIN, Value: true})
}

func (cam *Camera) ShowMainCmd() {
	cam.Command(CameraCmd{Action: HIDEMAIN, Value: false})
}

func (cam *Camera) Open() (err error) {
	cam.Decoder, err = mjpeg.NewDecoderFromURL(cam.URL)
	if err != nil {
		log.Println("NewDecoderFromURL", err)
	}
	return
}

func (cam *Camera) Close() {
	cam.MainHook.Close(0)
	cam.ThumbHook.Close(0)
	log.Printf("Closed '%s'\n", cam.Name)
}

const (
	delayNormal    = time.Millisecond * 1
	delayRetry     = time.Second
	delayHibernate = time.Second * 30
	recordLimit    = time.Minute
)

func (cam *Camera) stopRecording() {
	cam.Recording = false
}

func (cam *Camera) startRecording() {

	if cam.Recording {
		return
	}

}

func (cam *Camera) doCmd(cmd CameraCmd) {
	switch cmd.Action {
	case GET:
	case SET:
	case HIDEMAIN:
		b, _ := cmd.Value.(bool)
		cam.HideMain = b
	case HIDETHUMB:
		b, _ := cmd.Value.(bool)
		cam.HideThumb = b
	case HIDEALL:
		b, _ := cmd.Value.(bool)
		cam.HideAll = b
	case RECORD_START:
		cam.startRecording()
	case RECORD_STOP:
		cam.stopRecording()
	}
}

func (cam *Camera) Serve() {
	if cam.Busy {
		return
	}

	var (
		cmd CameraCmd
	)

	err := cam.Open()
	if err != nil {
		return
	}

	cam.Busy = true
	defer func() {
		cam.Busy = false
		cam.Close()
	}()

	var (
		delay = delayNormal
	)

	for {
		time.Sleep(delay)

		select {
		case <-cam.Quit:
			return
		case cmd = <-cam.Cmd:
			cam.doCmd(cmd)
			continue
		default:
		}

		if cam.HideAll {
			continue
		}

		img, err := cam.Decoder.Decode()

		if err != nil {
			log.Println(err)
			continue
		}

		if cam.MainHook != nil && !cam.HideMain {
			cam.MainHook.Update(img)
		}
		if cam.ThumbHook != nil && !cam.HideThumb {
			cam.ThumbHook.Update(img)
		}
	}

}
