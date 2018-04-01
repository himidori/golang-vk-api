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

var (
	userFields = "nickname,screen_name,sex,bdate,city,country," +
		"photo,photo_medium,photo_big,has_mobile,contacts," +
		"education,online,relation,last_seen,activity," +
		"can_write_private_message,can_see_all_posts,can_post,universities"
)

type User struct {
	UID                     int          `json:"id"`
	FirstName               string       `json:"first_name"`
	LastName                string       `json:"last_name"`
	Sex                     int          `json:"sex"`
	Nickname                string       `json:"nickname"`
	ScreenName              string       `json:"screen_name"`
	BDate                   string       `json:"bdate"`
	City                    *UserCity    `json:"city"`
	Country                 *UserCountry `json:"country"`
	Photo                   string       `json:"photo"`
	PhotoMedium             string       `json:"photo_medium"`
	PhotoBig                string       `json:"photo_big"`
	HasMobile               int          `json:"has_mobile"`
	Online                  int          `json:"online"`
	CanPost                 int          `json:"can_post"`
	CanSeeAllPosts          int          `json:"can_see_all_posts"`
	CanWritePrivateMessages int          `json:"can_write_private_message"`
	Status                  string       `json:"activity"`
	LastOnline              *LastSeen    `json:"last_seen"`
	Hidden                  int          `json:"hidden"`
	Deactivated             string       `json:"deactivated"`
	Relation                int          `json:"relation"`
}

type UserCity struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type UserCountry struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type LastSeen struct {
	Time     int64 `json:"time"`
	Platform int   `json:"platform"`
}

func (client *VKClient) UsersGet(users []int) ([]*User, error) {
	idsString := ArrayToStr(users)
	v := url.Values{}
	v.Add("user_ids", idsString)
	v.Add("fields", userFields)

	resp, err := client.makeRequest("users.get", v)
	if err != nil {
		return nil, err
	}

	var userList []*User
	json.Unmarshal(resp.Response, &userList)

	return userList, nil
}
