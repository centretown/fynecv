package appdata

import (
	"fmt"
	"fynecv/comm"
	"fynecv/vision"
	"log"
	"time"
)

type AppData struct {
	Cameras  []*vision.Camera
	Lights   []*Light
	Actions  []*Number
	Entities map[string]*Entity[any]
	client   *comm.WebSockClient
	stop     chan int
}

func NewAppData() *AppData {
	var data = &AppData{
		Cameras: []*vision.Camera{
			vision.NewCamera("http://192.168.0.7:9000/"),
			vision.NewCamera("http://192.168.0.7:9000/1/"),
		},
		Lights: NewLightList([]string{
			"light.led_matrix_24",
			"light.led_strip_24"}),
		Actions: NewNumberList([]string{
			"number.pan",
			"number.tilt"}),

		Entities: make(map[string]*Entity[any]),

		stop: make(chan int),
	}

	var err error
	data.client, err = comm.NewWebSockClient()
	if err != nil {
		log.Println("NewAppData", err)
	}
	return data
}

func (data *AppData) Load() {

	data.Lights = NewLightList([]string{
		"light.led_matrix_24",
		"light.led_strip_24"})
	data.Actions = NewNumberList([]string{
		"number.pan",
		"number.tilt"})
}

const (
	auth      = `{ "type":"auth", "access_token":"%s" }`
	config    = `{ "type": "get_config","id":%d }`
	states    = `{ "type":"get_states", "id":%d }`
	subscribe = `{ "type":"subscribe_events", "event_type":"state_changed", "id":%d }`
)

func (data *AppData) Monitor() {
	hs, err := comm.NewWebSockClient()
	if err != nil {
		log.Fatal(err, "Dial")
	}

	data.client = hs

	cmd := fmt.Sprintf(auth, Token)
	// auth twice to get response
	hs.Write(cmd)
	hs.Read()
	hs.Write(cmd)
	hs.Read()
	// collect current states
	go data.monitor()
	hs.WriteID(states)
	// hs.Read()
	// subscribe to state changes
	hs.WriteID(subscribe)
	// hs.Read()
}

func (data *AppData) monitor() {
	for {
		time.Sleep(time.Millisecond)

		select {
		case <-data.stop:
			return
		default:
			buf, err := data.client.Read()
			if err != nil {
				log.Println(err)
				continue
			}
			Parse(buf)
		}
	}
}

func (data *AppData) StopMonitor() {
	data.stop <- 1
}
