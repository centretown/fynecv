package appdata

import (
	"encoding/json"
	"fmt"
	"fynecv/comm"
	"fynecv/vision"
	"log"
	"time"

	"fyne.io/fyne/v2/data/binding"
)

type AppData struct {
	Cameras      []*vision.Camera
	Entities     map[string]*Entity[json.RawMessage]
	sock         *comm.WebSockClient
	stop         chan int
	loadStatesID int
	eventsID     int
	loaded       binding.Bool
	// Ready        binding.Bool
	Err error

	subscriptions map[string][]*Subscription
}

func NewAppData() *AppData {
	var data = &AppData{
		Cameras: []*vision.Camera{
			vision.NewCamera("http://192.168.0.7:9000/"),
			vision.NewCamera("http://192.168.0.7:9000/1/"),
		},

		Entities: make(map[string]*Entity[json.RawMessage]),

		stop:          make(chan int),
		subscriptions: make(map[string][]*Subscription),

		loaded: binding.NewBool(),
		// Ready:  binding.NewBool(),
	}

	var err error
	data.sock, err = comm.NewWebSockClient()
	if err != nil {
		log.Println("NewAppData", err)
	}
	return data
}

func (data *AppData) Subscribe(entityID string, subscription *Subscription) {
	list, ok := data.subscriptions[entityID]
	if !ok {
		list = make([]*Subscription, 1)
		list[0] = subscription
	} else {
		list = append(list, subscription)
	}
	data.subscriptions[entityID] = list
}

func (data *AppData) Consume(entityID string, newState *Entity[json.RawMessage]) {
	subs, ok := data.subscriptions[entityID]
	if ok {
		for _, sub := range subs {
			sub.Consume(newState)
		}
	}
}

func (data *AppData) CallService(cmd string) {
	data.sock.WriteID(cmd)
}

const (
	auth      = `{ "type":"auth", "access_token":"%s" }`
	config    = `{ "type":"get_config", "id":%d }`
	states    = `{ "type":"get_states", "id":%d }`
	subscribe = `{ "type":"subscribe_events", "event_type":"state_changed", "id":%d }`
)

func (data *AppData) StopMonitor() {
	data.stop <- 1
}

func (data *AppData) Monitor() error {

	defer func() {
		if data.Err != nil {
			log.Println(data.Err, "Monitor")
		}
	}()

	data.sock, data.Err = comm.NewWebSockClient()
	if data.Err != nil {
		return data.Err
	}

	// authorize
	cmd := fmt.Sprintf(auth, comm.Token)
	data.Err = data.sock.Write(cmd)
	if data.Err != nil {
		return data.Err
	}

	_, data.Err = data.sock.Read()
	if data.Err != nil {
		return data.Err
	}
	data.Err = data.sock.Write(cmd)
	if data.Err != nil {
		return data.Err
	}
	_, data.Err = data.sock.Read()
	if data.Err != nil {
		return data.Err
	}

	go data.monitor()

	data.loadStatesID, _ = data.sock.WriteID(states)
	data.loaded.AddListener(binding.NewDataListener(func() {
		isLoaded, _ := data.loaded.Get()
		if isLoaded {
			data.eventsID, data.Err = data.sock.WriteID(subscribe)
			if data.Err != nil {
				return
			}
		}
	}))
	return data.Err
}

func (data *AppData) monitor() {
	var (
		errCount int
		delay    time.Duration = time.Millisecond
	)
	for {
		time.Sleep(delay)

		select {
		case <-data.stop:
			return
		default:
			buf, err := data.sock.Read()
			if err != nil {
				errCount++
				if errCount > 10 {
					log.Fatal(err)
				}
				log.Println(err)
				continue
			}
			errCount = 0
			data.ParseResponse(buf)
		}
	}
}
