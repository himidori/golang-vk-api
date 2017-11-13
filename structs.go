package vkapi

import (
	"encoding/json"
)

type APIResponse struct {
	Response json.RawMessage `json:"response"`
}

type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	UID              int    `json:"user_id"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
