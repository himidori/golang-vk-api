package vkapi

import (
	"encoding/json"
	"net/url"
)

//last seen device
const (
	_ = iota
	PlatformMobile
	PlatformIPhone
	PlatfromIPad
	PlatformAndroid
	PlatformWPhone
	PlatformWindows
	PlatformWeb
)

type User struct {
	UID            int       `json:"uid"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Sex            int       `json:"sex"`
	Nickname       string    `json:"nickname"`
	ScreenName     string    `json:"screen_name"`
	BDate          string    `json:"bdate"`
	City           int       `json:"city"`
	Country        int       `json:"country"`
	Photo          string    `json:"photo"`
	PhotoMedium    string    `json:"photo_medium"`
	PhotoBig       string    `json:"photo_big"`
	HasMobile      int       `json:"has_mobile"`
	Online         int       `json:"online"`
	CanPost        int       `json:"can_post"`
	CanSeeAllPosts int       `json:"can_see_all_posts"`
	Status         string    `json:"activity"`
	LastOnline     *LastSeen `json:"last_seen"`
	Hidden         int       `json:"hidden"`
}

type LastSeen struct {
	Time     int64 `json:"time"`
	Platform int   `json:"platform"`
}

func (client *VKClient) GetUsers(users string) ([]User, error) {
	fields := "nickname,screen_name,sex,bdate,city,country,photo,photo_medium,photo_big,has_mobile,contacts,education,online,relation,last_seen,activity,can_write_private_messages,can_see_all_posts,can_post,universities"
	v := url.Values{}
	v.Add("user_ids", users)
	v.Add("fields", fields)

	resp, err := client.makeRequest("users.get", v)
	if err != nil {
		return []User{}, err
	}

	var userList []User
	err = json.Unmarshal(resp.Response, &userList)
	if err != nil {
		return []User{}, err
	}

	return userList, nil
}
