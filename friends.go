package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type FriendsRequests struct {
	Count    int           `json:"count"`
	Requests []*FriendItem `json:"items"`
}

type FriendItem struct {
	UserID        int     `json:"user_id"`
	MutualFriends *Mutual `json:"mutual"`
}

type Mutual struct {
	Count int   `json:"count"`
	Users []int `json:"users"`
}

func (client *VKClient) GetFriendRequests(count int, out int) (*FriendsRequests, error) {
	params := url.Values{}
	params.Set("count", strconv.Itoa(count))
	params.Set("out", strconv.Itoa(out))
	params.Set("extended", "1")

	resp, err := client.makeRequest("friends.getRequests", params)
	if err != nil {
		return nil, err
	}

	var reqs *FriendsRequests
	json.Unmarshal(resp.Response, &reqs)
	return reqs, nil
}

func (client *VKClient) AcceptFriendRequest(userID int, text string, follow int) error {
	params := url.Values{}
	params.Set("user_id", strconv.Itoa(userID))
	params.Set("follow", strconv.Itoa(follow))
	if text != "" {
		params.Set("text", text)
	}

	_, err := client.makeRequest("friends.add", params)
	if err != nil {
		return err
	}

	return nil
}
