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
	Lights       []*Light
	Actions      []*Number
	Entities     map[string]*Entity[json.RawMessage]
	client       *comm.WebSockClient
	stop         chan int
	loadStatesID int
	eventsID     int
	loaded       binding.Bool
	Ready        binding.Bool
	Err          error

	subscriptions map[string][]*Subscription
}

func NewAppData() *AppData {
	var data = &AppData{
		Cameras: []*vision.Camera{
			vision.NewCamera("http://192.168.0.7:9000/"),
			vision.NewCamera("http://192.168.0.7:9000/1/"),
		},

		Lights:   make([]*Light, 0),
		Actions:  make([]*Number, 0),
		Entities: make(map[string]*Entity[json.RawMessage]),

		stop:          make(chan int),
		subscriptions: make(map[string][]*Subscription),

		loaded: binding.NewBool(),
		Ready:  binding.NewBool(),
	}

	var err error
	data.client, err = comm.NewWebSockClient()
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

func (data *AppData) CallService(cmd string) (int, error) {
	return data.client.WriteID(cmd)
}

func (data *AppData) getRaw(entityID string) (ent *Entity[json.RawMessage]) {
	var ok bool
	ent, ok = data.Entities[entityID]
	if !ok {
		data.Err = fmt.Errorf("%s not found", entityID)
		ent = &Entity[json.RawMessage]{}
	}
	return
}

func (data *AppData) LoadLists() {
	data.LoadLightList()
	ShowYaml(data.Lights)
	data.LoadNumberList()
	ShowYaml(data.Actions)

	weather := &Weather{}
	weather.Entity.Copy(data.getRaw("weather.forecast_home"))
	ShowYaml(weather)

	zone := &Zone{}
	zone.Entity.Copy(data.getRaw("zone.home"))
	ShowYaml(zone)
}

func (data *AppData) LoadLightList() {
	lights := make([]*Light, 0)
	lightIDs := []string{
		"light.led_matrix_24",
		"light.led_strip_24"}

	for _, id := range lightIDs {
		light := &Light{}
		entity, ok := data.Entities[id]
		if ok {
			light.Entity.Copy(entity)
			lights = append(lights, light)
		}
	}
	data.Lights = lights
}

func (data *AppData) LoadNumberList() {
	actions := make([]*Number, 0)
	numberIDs := []string{
		"number.pan",
		"number.tilt"}
	for _, id := range numberIDs {
		entity, ok := data.Entities[id]
		if ok {
			number := &Number{}
			number.Entity.Copy(entity)
			actions = append(actions, number)
		}
	}
	data.Actions = actions
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

func (data *AppData) Monitor() {
	hs, err := comm.NewWebSockClient()
	if err != nil {
		log.Fatal(err, "Dial")
	}

	data.client = hs

	cmd := fmt.Sprintf(auth, Token)
	// authorize
	hs.Write(cmd)
	buf, _ := hs.Read()
	log.Println(buf, "AUTH1")
	hs.Write(cmd)
	buf, _ = hs.Read()
	log.Println(buf, "AUTH 2")

	go data.monitor()
	data.loadStatesID, _ = hs.WriteID(states)
	data.loaded.AddListener(binding.NewDataListener(func() {
		isLoaded, _ := data.loaded.Get()
		if isLoaded {
			data.LoadLists()
			data.eventsID, _ = hs.WriteID(subscribe)
			data.Ready.Set(true)
		}
	}))
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
			buf, err := data.client.Read()
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
