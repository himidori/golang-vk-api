package vkapi

import (
	"net/url"
	"strconv"
)

func (client *VKClient) SendGroupInvite(groupID int, userID int) error {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("user_id", strconv.Itoa(userID))
	_, err := client.makeRequest("groups.invite", params)
	if err != nil {
		return err
	}

	return nil
}
