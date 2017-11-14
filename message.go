package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Message struct {
	MID               int                `json:"mid"`
	Date              int64              `json:"date"`
	Out               int                `json:"out"`
	UID               int                `json:"uid"`
	ReadState         int                `json:"read_state"`
	Title             string             `json:"title"`
	Body              string             `json:"body"`
	ChatID            int                `json:"chat_id"`
	ChatActive        string             `json:"chat_active"`
	PushSettings      *Push              `json:"push_settings"`
	UsersCount        int                `json:"users_count"`
	AdminID           int                `json:"admin_id"`
	Photo50           string             `json:"photo_50"`
	Photo100          string             `json:"photo_100"`
	Photo200          string             `json:"photo_200"`
	ForwardedMessages []ForwardedMessage `json:"fwd_messages"`
	Attachments       []Attachment       `json:"attachments"`
}

type Push struct {
	Sound         int   `json:"sound"`
	DisabledUntil int64 `json:"disabled_until"`
}

type ForwardedMessage struct {
	UID  int    `json:"uid"`
	Date int64  `json:"date"`
	Body string `json:"body"`
}

type Attachment struct {
	Type     string           `json:"type"`
	Audio    *AudioAttachment `json:"audio"`
	Video    *VideoAttachment `json:"video"`
	Photo    *PhotoAttachment `json:"photo"`
	Document *DocAttachment   `json:"doc"`
	Link     *LinkAttachment  `json:"link"`
}

type AudioAttachment struct {
	AudioID   int    `json:"aid"`
	OwnerID   int    `json:"owner_id"`
	Artist    string `json:"artist"`
	Title     string `json:"title"`
	Duration  int    `json:"duration"`
	URL       string `json:"url"`
	Performer string `json:"performer"`
}

type VideoAttachment struct {
	VideoID     int    `json:"vid"`
	OwnerID     int    `json:"owner_id"`
	Title       string `json:"title"`
	Duration    int    `json:"duration"`
	Description string `json:"description"`
	Date        int64  `json:"date"`
	Views       int    `json:"views"`
	Image       string `json:"image"`
	ImageBig    string `json:"image_big"`
	ImageSmall  string `json:"image_small"`
	ImageXBig   string `json:"image_xbig"`
	AccessKey   string `json:"access_key"`
	Platform    string `json:"platform"`
	CanEdit     int    `json:"can_edit"`
}

type PhotoAttachment struct {
	PhotoID     int    `json:"pid"`
	AID         int    `json:"aid"`
	OwnerID     int    `json:"owner_id"`
	Source      string `json:"src"`
	SourceBig   string `json:"src_big"`
	SourceSmall string `json:"src_small"`
	SourceXBig  string `json:"src_xbig"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Text        string `json:"text"`
	Created     int64  `json:"created"`
	AccessKey   string `json:"access_key"`
}

type DocAttachment struct {
	DocID      int    `json:"did"`
	OwnerID    int    `json:"owner_id"`
	Title      string `json:"title"`
	Size       int    `json:"size"`
	Extenstion string `json:"ext"`
	URL        string `json:"url"`
	Date       int64  `json:"date"`
	AccessKey  string `json:"access_key"`
}

type LinkAttachment struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Target      string `json:"target"`
}

func (client *VKClient) GetDialogs(count int, params url.Values) ([]Message, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Add("count", strconv.Itoa(count))

	resp, err := client.makeRequest("messages.getDialogs", params)
	if err != nil {
		return []Message{}, err
	}
	var dialogs []Message
	clearedString := deleteFirstKey(string(resp.Response))
	err = json.Unmarshal([]byte(clearedString), &dialogs)
	if err != nil {
		return []Message{}, err
	}

	return dialogs, nil
}

func (client *VKClient) GetMessages(count int, params url.Values) ([]Message, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Add("count", strconv.Itoa(count))

	resp, err := client.makeRequest("messages.get", params)
	if err != nil {
		return []Message{}, err
	}

	var messages []Message
	clearedString := deleteFirstKey(string(resp.Response))
	err = json.Unmarshal([]byte(clearedString), &messages)
	if err != nil {
		return []Message{}, err
	}

	return messages, nil
}
