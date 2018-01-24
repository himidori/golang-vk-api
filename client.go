package vkapi

import (
	"encoding/json"
	"errors"
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
		return nil, err
	}

	vkclient.Self = token

	me, err := vkclient.GetUsers(strconv.Itoa(vkclient.Self.UID))
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
		Client: &http.Client{},
	}

	res, err := vkclient.isTokenValid(token)

	if !res {
		return nil, err
	}
	vkclient.Self.AccessToken = token

	return vkclient, nil
}

func (client *VKClient) isTokenValid(token string) (bool, error) {
	req, err := http.NewRequest("GET", "https://api.vk.com/method/users.get", nil)
	if err != nil {
		return false, err
	}
	q := req.URL.Query()
	q.Add("access_token", token)
	q.Add("v", "5.71")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Client.Do(req)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var apiresp APIResponse
	json.Unmarshal(body, &apiresp)
	if apiresp.ResponseError.ErrorCode != 0 {
		return false, errors.New("auth error: " + apiresp.ResponseError.ErrorMsg)
	}

	var user []User
	json.Unmarshal(apiresp.Response, &user)
	client.Self.UID = user[0].UID

	return true, nil
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
	json.Unmarshal(body, &token)

	if token.Error != "" {
		return Token{}, errors.New(token.Error + ": " + token.ErrorDescription)
	}

	return token, nil
}

func (client *VKClient) makeRequest(method string, params url.Values) (APIResponse, error) {
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

func deleteFirstKey(s string) string {
	return "[" + s[strings.Index(s, ",")+1:len(s)-1] + "]"
}
