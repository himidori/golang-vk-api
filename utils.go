package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

type ResolveScreenName struct {
	Type     string `json:"type"`
	ObjectID int    `json:"object_id"`
}

func (client *VKClient) ResolveScreenName(name string) (ResolveScreenName, error) {
	var res ResolveScreenName
	params := url.Values{}
	params.Set("screen_name", name)

	resp, err := client.makeRequest("utils.resolveScreenName", params)
	if err == nil {

		json.Unmarshal(resp.Response, &res)
	}
	return res, err

}

func ArrayToStr(a []int) string {
	s := []string{}
	for _, num := range a {
		s = append(s, strconv.Itoa(num))
	}
	return strings.Join(s, ",")
}
