package vkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	defer resp.Body.Close()

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
