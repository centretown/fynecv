package appdata

import (
	"time"
)

type Response struct {
	ID   int    `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

type Result struct {
	Response
	Success bool           `json:"success" yaml:"success"`
	Result  []*Entity[any] `json:"result" yaml:"result"`
}

type DataEntityID struct {
	EntityID string `json:"entity_id" yaml:"entity_id"`
}

type DataStateChange[T any] struct {
	EntityID string     `json:"entity_id" yaml:"entity_id"`
	OldState *Entity[T] `json:"old_state" yaml:"old_state"`
	NewState *Entity[T] `json:"new_state" yaml:"new_state"`
}

type Context struct {
	ID       string `json:"id" yaml:"id"`
	ParentID string `json:"parent_id" yaml:"parent_id"`
	UserID   string `json:"user_id" yaml:"user_id"`
}

type Event[T any] struct {
	EventType string    `json:"event_type" yaml:"event_type"`
	Origin    string    `json:"origin" yaml:"origin"`
	TimeFired time.Time `json:"time_fired" yaml:"time_fired"`
	Context   Context   `json:"context" yaml:"context"`
	Data      T         `json:"data" yaml:"data"`
}

type EventResult[T any] struct {
	Response
	Event Event[T] `json:"event" yaml:"event"`
}

type LightEventResult struct {
	EventResult[LightAttributes]
}

type Entity[T any] struct {
	EntityID     string    `json:"entity_id" yaml:"entity_id"`
	State        string    `json:"state" yaml:"state"`
	LastChanged  time.Time `json:"last_changed" yaml:"last_changed"`
	LastReported time.Time `json:"last_reported" yaml:"last_reported"`
	LastUpdated  time.Time `json:"last_updated" yaml:"last_updated"`
	Context      Context   `json:"context" yaml:"context"`
	Attributes   T         `json:"attributes" yaml:"attributes"`
}

type LightAttributes struct {
	Name       string    `json:"friendly_name" yaml:"friendly_name"`
	Brightness int       `json:"brightness" yaml:"brightness"`
	ColorMode  string    `json:"color_mode" yaml:"color_mode"`
	Effect     string    `json:"effect" yaml:"effect"`
	EffectList []string  `json:"effect_list" yaml:"effect_list"`
	ColorRGB   []uint8   `json:"rgb_color" yaml:"rgb_color"`
	ColorXY    []float64 `json:"xy_color" yaml:"xy_color"`
	ColorHS    []float64 `json:"hs_color" yaml:"hs_color"`
}

type Light struct {
	Entity[LightAttributes]
}

type NumberAttributes struct {
	Min   float64 `json:"min" yaml:"min"`
	Max   float64 `json:"max" yaml:"max"`
	Step  float64 `json:"step" yaml:"step"`
	Mode  string  `json:"mode" yaml:"mode"`
	Units string  `json:"unit_of_measurement" yaml:"unit_of_measurement"`
	Name  string  `json:"friendly_name" yaml:"friendly_name"`
}

type Number struct {
	Entity[NumberAttributes]
}

type SensorAttributes struct {
	StateClass  string `json:"state_class" yaml:"state_class"`
	Units       string `json:"unit_of_measurement" yaml:"unit_of_measurement"`
	DeviceClass string `json:"device_class" yaml:"device_class"`
	Name        string `json:"friendly_name" yaml:"friendly_name"`
}
