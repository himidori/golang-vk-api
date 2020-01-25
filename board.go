package vkapi

import (
	"crypto/rand"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

type Topic struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Created      int64  `json:"created"`
	CreatedBy    int    `json:"created_by"`
	Updated      int64  `json:"updated"`
	UpdatedBy    int    `json:"updated_by"`
	IsClosed     int    `json:"is_closed"`
	IsFixed      int    `json:"is_fixed"`
	Comments     int    `json:"comments"`
	FirstComment string `json:"first_comment"`
	LastComment  string `json:"last_comment"`
}

type Topics struct {
	Count        int      `json:"count"`
	Topics       []*Topic `json:"items"`
	DefaultOrder int      `json:"default_order"`
	CanAddTopics int      `json:"can_add_topics"`
	Profiles     []*User  `json:"profiles"`
}

type TopicComment struct {
	ID         int    `json:"id"`
	FromID     int    `json:"from_id"`
	Date       int64  `json:"date"`
	Text       string `json:"text"`
	ReplyToUID int    `json:"reply_to_uid"`
	ReplyToCID int    `json:"reply_to_cid"`
	// Attachments []*Attachments `json:"attachments"`
}

type Comments struct {
	Count    int             `json:"count"`
	Comments []*TopicComment `json:"items"`
	// Poll     *Poll           `json:"poll"`
	Profiles []*User `json:"profiles"`
}

func (client *VKClient) BoardAddTopic(groupID int, title string, text string, fromGroup bool, attachments []string) (int, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("title", title)
	params.Set("text", text)
	params.Set("from_group", strconv.Itoa(BoolToInt(fromGroup)))
	params.Set("attachments", strings.Join(attachments, ","))

	resp, err := client.MakeRequest("board.addTopic", params)
	if err != nil {
		return 0, err
	}

	var id int
	if err = json.Unmarshal(resp.Response, &id); err != nil {
		return 0, err
	}

	return id, nil
}

func (client *VKClient) BoardCloseTopic(groupID int, topicID int) (bool, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))

	resp, err := client.MakeRequest("board.closeTopic", params)
	if err != nil {
		return false, err
	}

	var ok int
	if err = json.Unmarshal(resp.Response, &ok); err != nil {
		return false, err
	}

	return IntToBool(ok), nil
}

func (client *VKClient) BoardCreateComment(groupID int, topicID int, message string, attachments []string, fromGroup bool, stickerID int) (int, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))
	params.Set("message", message)
	params.Set("from_group", strconv.Itoa(BoolToInt(fromGroup)))
	params.Set("sticker_id", strconv.Itoa(stickerID))
	params.Set("attachments", strings.Join(attachments, ","))

	guid := make([]byte, 16)
	_, err := rand.Read(guid)
	if err != nil {
		return 0, err
	}

	params.Set("guid", string(guid))

	resp, err := client.MakeRequest("board.createComment", params)
	if err != nil {
		return 0, err
	}

	var id int
	if err = json.Unmarshal(resp.Response, &id); err != nil {
		return 0, err
	}

	return id, nil
}

func (client *VKClient) BoardDeleteComment(groupID int, topicID int, commetID int) (bool, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))
	params.Set("comment_id", strconv.Itoa(commetID))

	resp, err := client.MakeRequest("board.deleteComment", params)
	if err != nil {
		return false, err
	}

	var ok int
	if err = json.Unmarshal(resp.Response, &ok); err != nil {
		return false, err
	}

	return IntToBool(ok), nil
}

func (client *VKClient) BoardDeleteTopic(groupID int, topicID int) (bool, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))

	resp, err := client.MakeRequest("board.deleteTopic", params)
	if err != nil {
		return false, err
	}

	var ok int
	if err = json.Unmarshal(resp.Response, &ok); err != nil {
		return false, err
	}

	return IntToBool(ok), nil
}

func (client *VKClient) BoardEditComment(groupID int, topicID int, commentID int, message string, attachments []string) (bool, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))
	params.Set("comment_id", strconv.Itoa(commentID))
	params.Set("message", message)
	params.Set("attachments", strings.Join(attachments, ","))

	resp, err := client.MakeRequest("board.editComment", params)
	if err != nil {
		return false, err
	}

	var ok int
	if err = json.Unmarshal(resp.Response, &ok); err != nil {
		return false, err
	}

	return IntToBool(ok), nil
}

func (client *VKClient) BoardEditTopic(groupID int, topicID int, title string) (bool, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))
	params.Set("title", title)

	resp, err := client.MakeRequest("board.editTopic", params)
	if err != nil {
		return false, err
	}

	var ok int
	if err = json.Unmarshal(resp.Response, &ok); err != nil {
		return false, err
	}

	return IntToBool(ok), nil
}

func (client *VKClient) BoardFixTopic(groupID int, topicID int) (bool, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))

	resp, err := client.MakeRequest("board.fixTopic", params)
	if err != nil {
		return false, err
	}

	var ok int
	if err = json.Unmarshal(resp.Response, &ok); err != nil {
		return false, err
	}

	return IntToBool(ok), nil
}

func (client *VKClient) BoardGetComments(groupID int, topicID int, count int, params url.Values) (*Comments, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))
	params.Set("count", strconv.Itoa(count))

	resp, err := client.MakeRequest("board.getComments", params)
	if err != nil {
		return nil, err
	}

	var comments *Comments
	json.Unmarshal(resp.Response, &comments)
	return comments, nil
}

func (client *VKClient) BoardGetTopics(groupID int, count int, params url.Values) (*Topics, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("count", strconv.Itoa(count))

	resp, err := client.MakeRequest("board.getTopics", params)
	if err != nil {
		return nil, err
	}

	var topics *Topics
	json.Unmarshal(resp.Response, &topics)
	return topics, nil
}

func (client *VKClient) BoardOpenTopic(groupID int, topicID int) (bool, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))

	resp, err := client.MakeRequest("board.openTopic", params)
	if err != nil {
		return false, err
	}

	var ok int
	if err = json.Unmarshal(resp.Response, &ok); err != nil {
		return false, err
	}

	return IntToBool(ok), nil
}

func (client *VKClient) BoardRestoreComment(groupID int, topicID int, commentID int) (bool, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))
	params.Set("comment_id", strconv.Itoa(commentID))

	resp, err := client.MakeRequest("board.restoreComment", params)
	if err != nil {
		return false, err
	}

	var ok int
	if err = json.Unmarshal(resp.Response, &ok); err != nil {
		return false, err
	}

	return IntToBool(ok), nil
}

func (client *VKClient) BoardUnfixTopic(groupID int, topicID int) (bool, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupID))
	params.Set("topic_id", strconv.Itoa(topicID))

	resp, err := client.MakeRequest("board.unfixTopic", params)
	if err != nil {
		return false, err
	}

	var ok int
	if err = json.Unmarshal(resp.Response, &ok); err != nil {
		return false, err
	}

	return IntToBool(ok), nil
}
