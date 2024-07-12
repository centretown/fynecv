package appdata

import (
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

func (data *AppData) ParseResponse(buf []byte) {
	resp := &Response{}
	data.Err = parseAny(buf, resp)
	if data.Err != nil {
		return
	}
	switch resp.Type {
	case "event":
		data.parseEvent(buf)
	case "result":
		data.parseResult(buf)
	}
}

func (data *AppData) parseResult(buf []byte) {
	result := &Result{}
	data.Err = parseAny(buf, result)
	if data.Err != nil {
		if result.ID == data.loadStatesID {
			data.loaded.Set(true)
		}
		return
	}

	if result.ID == data.loadStatesID {
		result := &StateResult{}
		data.Err = parseAny(buf, result)
		if data.Err != nil {
			log.Println(data.Err, "parseResult")
			return
		}

		for _, entity := range result.Entities {
			data.Entities[entity.EntityID] = entity
			data.Consume(entity.EntityID, entity)
		}
		data.loaded.Set(true)
	} else {
		ShowYaml(result)
	}
}

func (data *AppData) parseEvent(buf []byte) {
	idResult := &EventResult[DataState]{}
	data.Err = parseAny(buf, idResult)
	if data.Err != nil {
		return
	}

	if idResult.Event.EventType == "state_changed" {
		entityID := idResult.Event.Data.EntityID
		result := &EventResult[DataStateChange[json.RawMessage]]{}
		data.Err = parseAny(buf, result)
		if data.Err != nil {
			return
		}

		newState := result.Event.Data.NewState
		data.Entities[entityID] = newState
		data.Consume(entityID, newState)
	}
}

func parseAny(buf []byte, lresult any) (err error) {
	err = json.Unmarshal(buf, lresult)
	if err != nil {
		log.Println(string(buf))
		log.Println("parseAny Unmarshal", err)
	}
	return
}

func ShowYaml(entity any) {
	out, err := yaml.Marshal(entity)
	if err != nil {
		log.Println("Marshal yaml", err)
		return
	}
	fmt.Println(string(out))
}
