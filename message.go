package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Dialog struct {
	Count    int    `json:"count"`
	Messages []Item `json:"items"`
}

type Message struct {
	Count    int             `json:"count"`
	Messages []DialogMessage `json:"items"`
}

type Item struct {
	Message *DialogMessage `json:"message"`
	InRead  int            `json:"in_read"`
	OutRead int            `json:"out_read"`
}

type DialogMessage struct {
	MID               int                 `json:"id"`
	Date              int64               `json:"date"`
	Out               int                 `json:"out"`
	UID               int                 `json:"user_id"`
	ReadState         int                 `json:"read_state"`
	Title             string              `json:"title"`
	Body              string              `json:"body"`
	RandomID          int                 `json:"random_id"`
	ChatID            int                 `json:"chat_id"`
	ChatActive        string              `json:"chat_active"`
	PushSettings      *Push               `json:"push_settings"`
	UsersCount        int                 `json:"users_count"`
	AdminID           int                 `json:"admin_id"`
	Photo50           string              `json:"photo_50"`
	Photo100          string              `json:"photo_100"`
	Photo200          string              `json:"photo_200"`
	ForwardedMessages []ForwardedMessage  `json:"fwd_messages"`
	Attachments       []MessageAttachment `json:"attachments"`
}

type Push struct {
	Sound         int   `json:"sound"`
	DisabledUntil int64 `json:"disabled_until"`
}

type ForwardedMessage struct {
	UID  int    `json:"user_id"`
	Date int64  `json:"date"`
	Body string `json:"body"`
}

type MessageAttachment struct {
	Type     string           `json:"type"`
	Audio    *AudioAttachment `json:"audio"`
	Video    *VideoAttachment `json:"video"`
	Photo    *PhotoAttachment `json:"photo"`
	Document *DocAttachment   `json:"doc"`
	Link     *LinkAttachment  `json:"link"`
}

type HistoryAttachment struct {
	Attachments []HistoryAttachmentItem `json:"items"`
	NextFrom    string                  `json:"next_from"`
}

type HistoryAttachmentItem struct {
	MID        int                `json:"message_id"`
	Attachment *MessageAttachment `json:"attachment"`
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
	PhotoID   int    `json:"pid"`
	AID       int    `json:"aid"`
	OwnerID   int    `json:"owner_id"`
	Photo75   string `json:"photo_75"`
	Photo130  string `json:"photo_130"`
	Photo604  string `json:"photo_604"`
	Photo807  string `json:"photo_807"`
	Photo1280 string `json:"photo_1280"`
	Photo2560 string `json:"photo_2560"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Text      string `json:"text"`
	Created   int64  `json:"created"`
	AccessKey string `json:"access_key"`
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

func (client *VKClient) GetDialogs(count int, params url.Values) (Dialog, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Add("count", strconv.Itoa(count))

	resp, err := client.makeRequest("messages.getDialogs", params)
	if err != nil {
		return Dialog{}, err
	}

	var dialog Dialog
	json.Unmarshal(resp.Response, &dialog)

	return dialog, nil
}

func (client *VKClient) GetHistoryAttachments(peerID int, mediaType string, count int, params url.Values) (HistoryAttachment, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Add("count", strconv.Itoa(count))
	params.Add("media_type", mediaType)
	params.Add("peer_id", strconv.Itoa(peerID))

	resp, err := client.makeRequest("messages.getHistoryAttachments", params)
	if err != nil {
		return HistoryAttachment{}, err
	}

	var att HistoryAttachment
	json.Unmarshal(resp.Response, &att)
	return att, nil
}

func (client *VKClient) GetMessages(count int, params url.Values) (Message, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Add("count", strconv.Itoa(count))

	resp, err := client.makeRequest("messages.get", params)
	if err != nil {
		return Message{}, err
	}

	var messages Message
	json.Unmarshal(resp.Response, &messages)

	return messages, nil
}

func (client *VKClient) SendMessage(user interface{}, message string, params url.Values) error {
	if params == nil {
		params = url.Values{}
	}
	params.Add("message", message)

	switch user.(type) {
	case int:
		params.Add("user_id", strconv.Itoa(user.(int)))
	case string:
		params.Add("domain", user.(string))
	}

	_, err := client.makeRequest("messages.send", params)
	if err != nil {
		return err
	}

	return nil
}
