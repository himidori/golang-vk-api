package vkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	tokenURL = "https://oauth.vk.com/token"
	apiURL   = "https://api.vk.com/method/%s"
)

type VKClient struct {
	Self   Token
	Client *http.Client
}

func NewVKClient(user string, password string) (*VKClient, error) {
	vkclient := &VKClient{
		Client: &http.Client{},
	}

	token, err := vkclient.auth(user, password)
	if err != nil {
		return &VKClient{}, err
	}

	vkclient.Self = token
	return vkclient, nil
}

func (client *VKClient) auth(user string, password string) (Token, error) {
	req, err := http.NewRequest("GET", tokenURL, nil)
	if err != nil {
		return Token{}, err
	}

	q := req.URL.Query()
	q.Add("grant_type", "password")
	q.Add("client_id", "3140623")
	q.Add("client_secret", "VeWdmVclDCtn6ihuP1nt")
	q.Add("username", user)
	q.Add("password", password)
	q.Add("v", "5.40")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Client.Do(req)
	if err != nil {
		return Token{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Token{}, err
	}

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return Token{}, err
	}
	return token, nil
}

func (client *VKClient) makeRequest(method string, params url.Values) (APIResponse, error) {
	endpoint := fmt.Sprintf(apiURL, method)
	if params == nil {
		params = url.Values{}
	}

	params.Add("access_token", client.Self.AccessToken)

	resp, err := client.Client.PostForm(endpoint, params)
	if err != nil {
		return APIResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	var apiresp APIResponse
	err = json.Unmarshal(body, &apiresp)
	if err != nil {
		return APIResponse{}, err
	}
	return apiresp, nil
}

func deleteFirstKey(s string) string {
	return "[" + s[strings.Index(s, ",")+1:len(s)-1] + "]"
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

func (client *VKClient) GetDialogs(offset int, count int, startMID int, onlyUnread bool) ([]Message, error) {
	v := url.Values{}

	if offset != 0 {
		v.Add("offset", strconv.Itoa(offset))
	}
	if count != 0 {
		v.Add("count", strconv.Itoa(count))
	}
	if startMID != 0 {
		v.Add("start_message_id", strconv.Itoa(startMID))
	}
	if onlyUnread {
		v.Add("unread", "1")
	}

	resp, err := client.makeRequest("messages.getDialogs", v)
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

func (client *VKClient) GetMessages(offset int, count int, timeOffset int, filters int, lastMID int) ([]Message, error) {
	v := url.Values{}

	if offset != 0 {
		v.Add("offset", strconv.Itoa(offset))
	}
	if count != 0 {
		v.Add("count", strconv.Itoa(count))
	}
	if timeOffset != 0 {
		v.Add("time_offset", strconv.Itoa(timeOffset))
	}
	if filters != 0 {
		v.Add("filters", strconv.Itoa(filters))
	}
	if lastMID != 0 {
		v.Add("last_message_id", strconv.Itoa(lastMID))
	}

	resp, err := client.makeRequest("messages.get", v)
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

func (client *VKClient) GetWallPosts(ownerID int, domain string, offset int, count int, filter string) ([]WallPost, error) {
	v := url.Values{}

	if ownerID != 0 {
		v.Add("owner_id", strconv.Itoa(ownerID))
	}
	if domain != "" {
		v.Add("domain", domain)
	}
	if offset != 0 {
		v.Add("offset", strconv.Itoa(offset))
	}
	if count != 0 {
		v.Add("count", strconv.Itoa(count))
	}
	if filter != "" {
		v.Add("filter", filter)
	}

	resp, err := client.makeRequest("wall.get", v)
	if err != nil {
		return []WallPost{}, err
	}

	var posts []WallPost
	clearedString := deleteFirstKey(string(resp.Response))
	err = json.Unmarshal([]byte(clearedString), &posts)
	if err != nil {
		return []WallPost{}, err
	}

	return posts, nil
}

func (client *VKClient) getLongPollServer() (LongPollServer, error) {
	resp, err := client.makeRequest("messages.getLongPollServer", nil)
	if err != nil {
		return LongPollServer{}, err
	}
	var server LongPollServer
	err = json.Unmarshal(resp.Response, &server)
	if err != nil {
		return LongPollServer{}, err
	}

	return server, nil
}

func (client *VKClient) ListenLongPollServer() (LongPollChannel, error) {
	ch := make(chan LongPollMessage, 10)
	server, err := client.getLongPollServer()
	if err != nil {
		return ch, err
	}

	go func() {
		for {
			req, err := http.NewRequest("GET", "https://"+server.Server, nil)
			if err != nil {
				return
			}

			q := req.URL.Query()
			q.Add("act", "a_check")
			q.Add("key", server.Key)
			q.Add("ts", strconv.FormatInt(server.TS, 10))
			q.Add("wait", "25")
			q.Add("mode", "2")
			q.Add("version", "1")
			req.URL.RawQuery = q.Encode()

			resp, err := client.Client.Do(req)
			if err != nil {
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}

			var updates LongPollUpdate
			err = json.Unmarshal(body, &updates)

			switch updates.Failed {
			case 0:
				for _, update := range updates.Updates {
					updateID := update[0].(float64)

					switch updateID {
					case 4: //new message
						var message LongPollMessage
						message.MessageID = int(update[1].(float64))
						message.UserID = int(update[3].(float64))
						message.Date = int64(update[4].(float64))
						message.Title = update[5].(string)
						message.Body = update[6].(string)
						message.Attachments = make(map[string]string)

						for k, v := range update[7].(map[string]interface{}) {
							message.Attachments[k] = v.(string)
						}

						ch <- message
					}
				}
				server.TS = updates.TS
			case 1:
				server.TS = updates.TS
			}
		}
	}()

	return ch, nil
}
