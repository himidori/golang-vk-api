package vkapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	tokenURL = "https://oauth.vk.com/token"
	apiURL   = "https://api.vk.com/method/%s"
)

const (
	DeviceIPhone = iota
	DeviceWPhone
	DeviceAndroid
)

type ratelimiter struct {
	requestsCount   int
	lastRequestTime time.Time
}

type VKClient struct {
	Self   Token
	Client *http.Client
	rl     *ratelimiter
	cb     *callbackHandler
}

func NewVKClient(device int, user string, password string) (*VKClient, error) {
	vkclient := &VKClient{
		Client: &http.Client{Timeout: time.Second * 10},
	}

	token, err := vkclient.auth(device, user, password)
	if err != nil {
		return nil, err
	}

	vkclient.Self = token

	vkclient.rl = &ratelimiter{}
	vkclient.cb = &callbackHandler{
		events: make(map[string]func(*LongPollMessage)),
	}

	me, err := vkclient.UsersGet([]int{vkclient.Self.UID})
	if err != nil {
		return nil, err
	}

	vkclient.Self.FirstName = me[0].FirstName
	vkclient.Self.LastName = me[0].LastName
	vkclient.Self.PicSmall = me[0].Photo
	vkclient.Self.PicMedium = me[0].PhotoMedium
	vkclient.Self.PicBig = me[0].PhotoBig

	return vkclient, nil
}

func NewVKClientWithToken(token string) (*VKClient, error) {
	vkclient := &VKClient{
		Client: &http.Client{Timeout: time.Second * 10},
	}

	err := vkclient.isTokenValid(token)
	if err != nil {
		return nil, err
	}

	vkclient.Self.AccessToken = token
	vkclient.rl = &ratelimiter{}
	vkclient.cb = &callbackHandler{
		events: make(map[string]func(*LongPollMessage)),
	}

	return vkclient, nil
}

func (client *VKClient) isTokenValid(token string) error {
	req, err := http.NewRequest("GET", "https://api.vk.com/method/users.get", nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("access_token", token)
	q.Add("v", "5.71")
	q.Add("fields", "photo,photo_medium,photo_big")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiresp APIResponse
	json.Unmarshal(body, &apiresp)
	if apiresp.ResponseError.ErrorCode != 0 {
		return errors.New("auth error: " + apiresp.ResponseError.ErrorMsg)
	}

	var user []User
	json.Unmarshal(apiresp.Response, &user)
	client.Self.UID = user[0].UID
	client.Self.FirstName = user[0].FirstName
	client.Self.LastName = user[0].LastName
	client.Self.PicSmall = user[0].Photo
	client.Self.PicMedium = user[0].PhotoMedium
	client.Self.PicBig = user[0].PhotoBig

	return nil
}

func (client *VKClient) auth(device int, user string, password string) (Token, error) {
	req, err := http.NewRequest("GET", tokenURL, nil)
	if err != nil {
		return Token{}, err
	}

	clientID := ""
	clientSecret := ""

	switch device {
	case DeviceIPhone:
		clientID = "3140623"
		clientSecret = "VeWdmVclDCtn6ihuP1nt"
	case DeviceWPhone:
		clientID = "3697615"
		clientSecret = "AlVXZFMUqyrnABp8ncuU"
	case DeviceAndroid:
		clientID = "2274003"
		clientSecret = "hHbZxrka2uZ6jB1inYsH"
	default:
		clientID = "3140623"
		clientSecret = "VeWdmVclDCtn6ihuP1nt"
	}

	q := req.URL.Query()
	q.Add("grant_type", "password")
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
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
	json.Unmarshal(body, &token)

	if token.Error != "" {
		return Token{}, errors.New(token.Error + ": " + token.ErrorDescription)
	}

	return token, nil
}

func (client *VKClient) makeRequest(method string, params url.Values) (APIResponse, error) {
	if client.rl.requestsCount == 3 {
		secs := time.Since(client.rl.lastRequestTime).Seconds()
		ms := int((1 - secs) * 1000)
		if ms > 0 {
			duration := time.Duration(ms * int(time.Millisecond))
			//fmt.Println("attempted to make more than 3 requests per second, "+
			//"sleeping for", ms, "ms")
			time.Sleep(duration)
		}

		client.rl.requestsCount = 0
	}

	endpoint := fmt.Sprintf(apiURL, method)
	if params == nil {
		params = url.Values{}
	}

	params.Set("access_token", client.Self.AccessToken)
	params.Set("v", "5.71")

	resp, err := client.Client.PostForm(endpoint, params)
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	client.rl.requestsCount++
	client.rl.lastRequestTime = time.Now()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	var apiresp APIResponse
	json.Unmarshal(body, &apiresp)

	if apiresp.ResponseError.ErrorCode != 0 {
		return APIResponse{}, errors.New("Error code: " + strconv.Itoa(apiresp.ResponseError.ErrorCode) + ", " + apiresp.ResponseError.ErrorMsg)
	}
	return apiresp, nil
}
