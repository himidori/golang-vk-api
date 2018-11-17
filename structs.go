package vkapi

import (
	"encoding/json"
)

type APIResponse struct {
	Response      json.RawMessage `json:"response"`
	ResponseError Error           `json:"error"`
}

type Error struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	UID              int    `json:"user_id"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	IsSerivceToken   bool
	FirstName        string
	LastName         string
	PicSmall         string
	PicMedium        string
	PicBig           string
}
