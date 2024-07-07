package appdata

import (
	"encoding/json"
	"fynecv/comm"
	"log"
)

func NewLightList(entities []string) []*Light {
	var lights = make([]*Light, 0, len(entities))
	for _, ent := range entities {

		buf, err := comm.Get("states/" + ent)
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
		buf, err := comm.Get("states/" + light.EntityID)
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
