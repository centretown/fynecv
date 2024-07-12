package appdata

import "encoding/json"

type ServiceCommand struct {
	Id      int    `json:"id" yaml:"id"`
	Type    string `json:"type" yaml:"type"`
	Domain  string `json:"domain" yaml:"domain"`
	Service string `json:"service" yaml:"service"`
	Target  struct {
		EntityID string `json:"entity_id" yaml:"entity_id"`
	}
	ServiceData    json.RawMessage `json:"service_data" yaml:"service_data"`
	ReturnResponse bool            `json:"return_response" yaml:"return_response"`
}
