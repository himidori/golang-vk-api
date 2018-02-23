package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type FriendsRequests struct {
	Count    int        `json:"count"`
	Requests []*Request `json:"items"`
}

type Friends struct {
	Count int     `json:"count"`
	Users []*User `json:"items"`
}

type Request struct {
	UserID        int     `json:"user_id"`
	MutualFriends *Mutual `json:"mutual"`
}

type Mutual struct {
	Count int   `json:"count"`
	Users []int `json:"users"`
}

func (client *VKClient) FriendsGet(uid int, count int) (int, []*User, error) {
	params := url.Values{}
	params.Set("count", strconv.Itoa(count))
	params.Set("fields", userFields)

	resp, err := client.makeRequest("friends.get", params)
	if err != nil {
		return 0, nil, err
	}

	var friends *Friends
	json.Unmarshal(resp.Response, &friends)
	return friends.Count, friends.Users, nil
}

func (client *VKClient) FriendsGetRequests(count int, out int) (int, []*Request, error) {
	params := url.Values{}
	params.Set("count", strconv.Itoa(count))
	params.Set("out", strconv.Itoa(out))
	params.Set("extended", "1")

	resp, err := client.makeRequest("friends.getRequests", params)
	if err != nil {
		return 0, nil, err
	}

	var reqs *FriendsRequests
	json.Unmarshal(resp.Response, &reqs)
	return reqs.Count, reqs.Requests, nil
}

func (client *VKClient) FriendsAdd(userID int, text string, follow int) error {
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

func (client *VKClient) FriendsDelete(userID int) error {
	params := url.Values{}
	params.Set("user_id", strconv.Itoa(userID))

	_, err := client.makeRequest("friends.delete", params)
	if err != nil {
		return err
	}

	return nil
}
