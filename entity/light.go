package entity

import (
	"encoding/json"
	"fynecv/hass"
	"log"
	"time"
)

type Light struct {
	EntityID   string `json:"entity_id"`
	State      string `json:"state"`
	Attributes struct {
		Name       string    `json:"friendly_name"`
		Brightness float64   `json:"brightness"`
		ColorMode  string    `json:"color_mode"`
		Effect     string    `json:"effect"`
		EffectList []string  `json:"effect_list"`
		ColorRGB   []uint8   `json:"rgb_color"`
		ColorXY    []float64 `json:"xy_color"`
		ColorHS    []float64 `json:"hs_color"`
	}
	LastChanged  time.Time `json:"last_changed"`
	LastReported time.Time `json:"last_reported"`
	LastUpdated  time.Time `json:"last_updated"`
}

func BuildLightList(entities []string) []*Light {

	var lights = make([]*Light, 0, len(entities))

	for _, ent := range entities {

		log.Println(ent)

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

// {"entity_id": "light.led_matrix_24","effect": "rainbow-vertical","brightness_pct": 50}
// const JS = `
// {
//   "attributes": {
//     "brightness": 110,
//     "color_mode": "rgb",
//     "effect": "None",
//     "effect_list": [
//       "None",
//       "black-white-scan",
//       "complementary-scan",
//       "double-scan",
//       "gradient-scan",
//       "rainbow-diagonal",
//       "rainbow-horizontal",
//       "rainbow-vertical",
//       "split_in_three",
//       "split_in_two",
//       "spotlight"
//     ],
//     "friendly_name": "LED Strip 24",
//     "hs_color": [
//       30.496,
//       94.902
//     ],
//     "rgb_color": [
//       255,
//       136,
//       13
//     ],
//     "supported_color_modes": [
//       "rgb"
//     ],
//     "supported_features": 44,
//     "xy_color": [
//       0.599,
//       0.382
//     ]
//   },
//   "context": {
//     "id": "01J1DS1MNNP9GC9TVAPNKCBA82",
//     "parent_id": null,
//     "user_id": "4fcc5ee6683d4c9eb1ac5e8e1e42d240"
//   },
//   "entity_id": "light.led_strip_24",
//   "last_changed": "2024-06-23T21:17:17.781514+00:00",
//   "last_reported": "2024-06-27T21:18:56.311986+00:00",
//   "last_updated": "2024-06-27T21:18:56.311986+00:00",
//   "state": "on"
// }
// `
