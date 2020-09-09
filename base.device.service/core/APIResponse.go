package core

import (
	"encoding/json"
	"log"
)

//APIResponse struct for api call response
type APIResponse struct {
	CallbackID string `json:"CallbackId"`
	Success    bool   `json:"Success"`
	Code       string `json:"Code"`
}

//Marshal APIResponse struct
func (data *APIResponse) Marshal() string {
	jMessage, jerr := json.Marshal(data)
	if jerr != nil {
		log.Println("Marshal APIResponse data error:", jerr)
		return ""
	}
	return string(jMessage)
}
