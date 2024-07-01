package entity

import (
	"encoding/json"
	"fynecv/hass"
	"log"
)

type Light struct {
	State
	Attributes struct {
		Name       string    `json:"friendly_name" yaml:"friendly_name"`
		Brightness float64   `json:"brightness" yaml:"brightness"`
		ColorMode  string    `json:"color_mode" yaml:"color_mode"`
		Effect     string    `json:"effect" yaml:"effect"`
		EffectList []string  `json:"effect_list" yaml:"effect_list"`
		ColorRGB   []uint8   `json:"rgb_color" yaml:"rgb_color"`
		ColorXY    []float64 `json:"xy_color" yaml:"xy_color"`
		ColorHS    []float64 `json:"hs_color" yaml:"hs_color"`
	}
	Context Context
}

func NewLightList(entities []string) []*Light {
	var lights = make([]*Light, 0, len(entities))
	for _, ent := range entities {

		buf, err := hass.Get("states/" + ent)
		if err != nil {
			log.Println(err, string(buf))
			continue
		}

		light := &Light{}
		err = json.Unmarshal(buf, light)
		if err != nil {
			log.Println(err, string(buf))
			continue
		}

		lights = append(lights, light)
	}
	return lights
}

func RefreshLights(lights []*Light) {
	for _, light := range lights {
		buf, err := hass.Get("states/" + light.EntityID)
		if err != nil {
			log.Println(err, string(buf))
			continue
		}

		err = json.Unmarshal(buf, light)
		if err != nil {
			log.Println(err, string(buf))
			continue
		}
	}
}
