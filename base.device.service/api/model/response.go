package model

import "time"

//Response struct for responses on API requests
type Response struct {
	Code            string    `json:"Code" example:"Device is offline"`
	Comment         string    `json:"Comment" example:"Response description"`
	ExecutedCommand string    `json:"ExecutedCommand" example:"TRACK,0#"`
	CallbackID      string    `json:"CallbackId" example:"1a3bf03d930d46adba3ee9f20fc50508"`
	CreatedAt       time.Time `json:"CreatedAt" example:"2018-05-09T01:02:03Z"`
	Success         bool      `json:"Success" example:"true"`
}
