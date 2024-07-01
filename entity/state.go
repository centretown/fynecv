package entity

import "time"

type State struct {
	EntityID     string    `json:"entity_id"`
	State        string    `json:"state"`
	LastChanged  time.Time `json:"last_changed"`
	LastReported time.Time `json:"last_reported"`
	LastUpdated  time.Time `json:"last_updated"`
}
