package appdata

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parse(buf []byte) {

	resp := &Response{}
	err := json.Unmarshal(buf, resp)
	if err != nil {
		log.Println("Parse", err)
		log.Println(string(buf))
		return
	}

	log.Println("Id:", resp.ID, "Type:", resp.Type)
	switch resp.Type {
	case "event":
		parseEvent(buf)
	case "result":
		parseResult(buf)
	}
}

func parseResult(buf []byte) {
	result := &Result{}

	err := json.Unmarshal(buf, result)
	if err != nil {
		log.Println("Result", err)
		log.Println(string(buf))
		return
	}
	showYaml(result)
	// log.Println("Id:", result.Id, "Type:", result.Type, "Success:", result.Success)

}

func parseEvent(buf []byte) {
	// fmt.Println(string(buf))
	result := &EventResult[DataEntityID]{}
	err := json.Unmarshal(buf, result)

	if err != nil {
		log.Println("EventResult", err)
		log.Println(string(buf))
		return
	}

	dump := func() {
		log.Println("Id:", result.ID,
			"Type:", result.Type,
			"EventType", result.Event.EventType,
			"EntityID", result.Event.Data.EntityID,
			"Origin", result.Event.Origin,
			"TimeFired", result.Event.TimeFired.Local())
	}

	if result.Event.EventType != "state_changed" {
		dump()
		return
	}

	entityID := result.Event.Data.EntityID
	periodPos := strings.Index(entityID, ".")
	if periodPos == -1 {
		log.Println("NO PERIOD FOUND")
		dump()
		return
	}

	entityType := entityID[:periodPos]
	switch entityType {
	case "light":
		lresult := &EventResult[DataStateChange[LightAttributes]]{}
		parseAny(buf, lresult)
		light := lresult.Event.Data.NewState
		showYaml(light)
	case "number":
		lresult := &EventResult[DataStateChange[NumberAttributes]]{}
		parseAny(buf, lresult)
		num := lresult.Event.Data.NewState
		showYaml(num)
	case "sensor":
		lresult := &EventResult[DataStateChange[SensorAttributes]]{}
		parseAny(buf, lresult)
		num := lresult.Event.Data.NewState
		showYaml(num)
	}
}

func parseAny(buf []byte, lresult any) {
	err := json.Unmarshal(buf, lresult)
	if err != nil {
		log.Println("Unmarshal", err)
		log.Println(string(buf))
		return
	}
}

func showYaml(entity any) {
	out, err := yaml.Marshal(entity)
	if err != nil {
		log.Println("Marshal yaml", err)
		return
	}

	fmt.Println(string(out))

}
