package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Gift struct {
	Count int         `json:"count"`
	Gifts []*GiftItem `json:"items"`
}

type GiftItem struct {
	ID      int       `json:"id"`
	FromID  int       `json:"from_id"`
	Message string    `json:"message"`
	Date    int64     `json:"date"`
	Info    *GiftInfo `json:"gift"`
	Privacy int       `json:"privacy"`
}

type GiftInfo struct {
	ID       int    `json:"id"`
	Thumb256 string `json:"thumb_256"`
	Thumb48  string `json:"thumb_48"`
	Thumb96  string `json:"thumb_96"`
}

func (client *VKClient) GetGifts(id int, count int, offset int) (*Gift, error) {
	params := url.Values{}

	params.Set("user_id", strconv.Itoa(id))
	params.Set("count", strconv.Itoa(count))
	params.Set("offset", strconv.Itoa(offset))

	resp, err := client.MakeRequest("gifts.get", params)
	if err != nil {
		return nil, err
	}

	var gift *Gift
	json.Unmarshal(resp.Response, &gift)

	return gift, nil
}