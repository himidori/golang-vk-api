package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const (
	TypePost         = "post"
	TypeComment      = "comment"
	TypePhoto        = "photo"
	TypeAudio        = "audio"
	TypeVideo        = "video"
	TypeNote         = "note"
	TypePhotoComment = "photo_comment"
	TypeVideoComment = "video_comment"
	TypeTopicComment = "topic_comment"
	TypeSitepage     = "sitepage"
)

func (client *VKClient) GetLikes(itemType string, ownerID int, itemID int, count int, params url.Values) (int, []int, error) {
	if params == nil {
		params = url.Values{}
	}

	params.Add("type", itemType)
	params.Add("count", strconv.Itoa(count))
	params.Add("owner_id", strconv.Itoa(ownerID))
	params.Add("item_id", strconv.Itoa(itemID))

	resp, err := client.makeRequest("likes.getList", params)
	if err != nil {
		return 0, []int{}, err
	}

	var data struct {
		LikesCount int   `json:"count"`
		IDs        []int `json:"users"`
	}

	json.Unmarshal(resp.Response, &data)

	return data.LikesCount, data.IDs, nil
}
