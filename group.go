package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

var (
	groupFields = "description,age_limits,activity,can_create_topic," +
		"can_message,can_post,can_see_all_posts,contacts,has_photo," +
		"is_messages_blocked"
)

type Group struct {
	ID                int             `json:"id"`
	Name              string          `json:"name"`
	ScreenName        string          `json:"screen_name"`
	Description       string          `json:"description"`
	Activity          string          `json:"activity"`
	Contacts          []*GroupContact `json:"contacts"`
	IsClosed          int             `json:"is_closed"`
	Type              string          `json:"type"`
	IsAdmin           int             `json:"is_admin"`
	IsMember          int             `json:"is_member"`
	HasPhoto          int             `json:"has_photo"`
	IsMessagesBlocked int             `json:"is_messages_blocked"`
	Photo50           string          `json:"photo_50"`
	Photo100          string          `json:"photo_100"`
	Photo200          string          `json:"photo_200"`
	AgeLimit          int             `json:"age_limits"`
	CanCreateTopic    int             `json:"can_create_topic"`
	CanMessage        int             `json:"can_message"`
	CanPost           int             `json:"can_post"`
	CanSeeAllPosts    int             `json:"can_see_all_posts"`
}

type GroupSearchResult struct {
	Count  int      `json:"count"`
	Groups []*Group `json:"items"`
}

type GroupContact struct {
	UID         int    `json:"user_id"`
	Description string `json:"desc"`
}

type GroupMembers struct {
	Count   int     `json:"count"`
	Members []*User `json:"items"`
}

func (client *VKClient) GroupSendInvite(groupID int, userID int) error {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("user_id", strconv.Itoa(userID))
	_, err := client.MakeRequest("groups.invite", params)
	if err != nil {
		return err
	}

	return nil
}

func (client *VKClient) GroupSearch(query string, count int) (int, []*Group, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("count", strconv.Itoa(count))
	params.Set("fields", groupFields)
	resp, err := client.MakeRequest("groups.search", params)
	if err != nil {
		return 0, nil, err
	}

	var res *GroupSearchResult
	json.Unmarshal(resp.Response, &res)
	return res.Count, res.Groups, nil
}

func (client *VKClient) GroupGet(userID int, count int) (int, []*Group, error) {
	params := url.Values{}
	params.Set("user_id", strconv.Itoa(userID))
	params.Set("count", strconv.Itoa(count))
	params.Set("extended", "1")
	resp, err := client.MakeRequest("groups.get", params)
	if err != nil {
		return 0, nil, err
	}

	var res *GroupSearchResult
	json.Unmarshal(resp.Response, &res)
	return res.Count, res.Groups, nil
}

func (client *VKClient) GroupsGetByID(groupsID []int) ([]*Group, error) {
	params := url.Values{}
	params.Set("group_ids", ArrayToStr(groupsID))
	resp, err := client.MakeRequest("groups.getById", params)
	if err != nil {
		return nil, err
	}

	var groupsList []*Group
	json.Unmarshal(resp.Response, &groupsList)

	return groupsList, nil
}

func (client *VKClient) GroupGetMembers(group_id int, count int) (int, []*User, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(group_id))
	params.Set("count", strconv.Itoa(count))
	params.Set("fields", userFields)
	resp, err := client.MakeRequest("groups.getMembers", params)
	if err != nil {
		return 0, nil, err
	}

	var res *GroupMembers
	json.Unmarshal(resp.Response, &res)
	return res.Count, res.Members, nil
}
