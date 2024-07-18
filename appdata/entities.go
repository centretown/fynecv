package appdata

import (
	"encoding/json"
	"time"
)

type Response struct {
	ID   int    `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

type Result struct {
	Response
	Success bool `json:"success" yaml:"success"`
	// Entities []*Entity[json.RawMessage] `json:"result" yaml:"result"`
}

type StateResult struct {
	Response
	Success  bool                       `json:"success" yaml:"success"`
	Entities []*Entity[json.RawMessage] `json:"result" yaml:"result"`
}

type DataState struct {
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

type EntityStore struct {
	Entity[json.RawMessage]
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

type NumberAttributes struct {
	Min   float64 `json:"min" yaml:"min"`
	Max   float64 `json:"max" yaml:"max"`
	Step  float64 `json:"step" yaml:"step"`
	Mode  string  `json:"mode" yaml:"mode"`
	Units string  `json:"unit_of_measurement" yaml:"unit_of_measurement"`
	Name  string  `json:"friendly_name" yaml:"friendly_name"`
}

// weather.forecast_home
type WeatherAttributes struct {
	Attribution       string  `json:"attribution" yaml:"attribution"`
	CloudCoverage     float64 `json:"cloud_coverage" yaml:"cloud_coverage"`
	DewPoint          float64 `json:"dew_point" yaml:"dew_point"`
	Name              string  `json:"friendly_name" yaml:"friendly_name"`
	Humidity          float64 `json:"humidity" yaml:"humidity"`
	PrecipitationUnit string  `json:"precipitation_unit" yaml:"precipitation_unit"`
	Pressure          float64 `json:"pressure" yaml:"pressure"`
	PressureUnit      string  `json:"pressure_unit" yaml:"pressure_unit"`
	SupportedFeatures int     `json:"supported_features" yaml:"supported_features"`
	Temperature       float64 `json:"temperature" yaml:"temperature"`
	TemperatureUnit   string  `json:"temperature_unit" yaml:"temperature_unit"`
	VisibilityUnit    string  `json:"visibility_unit" yaml:"visibility_unit"`
	WindBearing       float64 `json:"wind_bearing" yaml:"wind_bearing"`
	WindSpeed         float64 `json:"wind_speed" yaml:"wind_speed"`
	WindSpeedUnit     string  `json:"wind_speed_unit" yaml:"wind_speed_unit"`
}

// sensor.wifi_signal_28

type SensorAttributes struct {
	StateClass  string `json:"state_class" yaml:"state_class"`
	Units       string `json:"unit_of_measurement" yaml:"unit_of_measurement"`
	DeviceClass string `json:"device_class" yaml:"device_class"`
	Name        string `json:"friendly_name" yaml:"friendly_name"`
}

type ZoneAttributes struct {
	Latitude  float64  `json:"latitude" yaml:"latitude"`
	Longitude float64  `json:"longitude" yaml:"longitude"`
	Radius    float64  `json:"radius" yaml:"radius"`
	Passive   bool     `json:"passive" yaml:"passive"`
	Persons   []string `json:"persons" yaml:"persons"`
	Editable  bool     `json:"editable" yaml:"editable"`
	Icon      string   `json:"icon" yaml:"icon"`
	Name      string   `json:"friendly_name" yaml:"friendly_name"`
}

type Light struct {
	Entity[LightAttributes]
}
type Number struct {
	Entity[NumberAttributes]
}
type Weather struct {
	Entity[WeatherAttributes]
}
type Wifi struct {
	Entity[SensorAttributes]
}
type AnyData struct {
	Entity[any]
}
type Zone struct {
	Entity[ZoneAttributes]
}

type Consumer interface {
	Copy(src *Entity[json.RawMessage])
}

func (dst *Entity[T]) Copy(src *Entity[json.RawMessage]) {
	dst.EntityID = src.EntityID
	dst.State = src.State
	dst.LastChanged = src.LastChanged
	dst.LastUpdated = src.LastUpdated
	dst.LastReported = src.LastReported
	dst.Context = src.Context
	json.Unmarshal(src.Attributes, &dst.Attributes)
}
