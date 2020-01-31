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
	DefaultOrder float64  `json:"default_order"`
	CanAddTopics int      `json:"can_add_topics"`
	Profiles     []*User  `json:"profiles"`
}

type TopicCommentLike struct {
	Count     int `json:"count"`
	UserLikes int `json:"user_likes"`
	CanLike   int `json:"can_like"`
}

type TopicComment struct {
	ID          int               `json:"id"`
	FromID      int               `json:"from_id"`
	Date        int64             `json:"date"`
	Text        string            `json:"text"`
	Likes       *TopicCommentLike `json:"likes"`
	ReplyToUID  int               `json:"reply_to_uid"`
	ReplyToCID  int               `json:"reply_to_cid"`
	Attachments []*Attachment     `json:"attachments"`
}

type Attachment struct {
	Type    string             `json:"type"`
	Photo   *AttachmentPhoto   `json:"photo"`
	Video   *AttachmentVideo   `json:"video"`
	Audio   *AttachmentAudio   `json:"audio"`
	Doc     *AttachmentDoc     `json:"doc"`
	Sticker *AttachmentSticker `json:"sticker"`
}

type AttachmentImageInfo struct {
	Type        string `json:"type"`
	Url         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	WithPadding int    `json:"with_padding"`
}

type AttachmentPhoto struct {
	ID        int                    `json:"id"`
	AlbumID   int                    `json:"album_id"`
	OwnerID   int                    `json:"owner_id"`
	UserID    int                    `json:"user_id"`
	Sizes     []*AttachmentImageInfo `json:"sizes"`
	Text      string                 `json:"text"`
	Date      int64                  `json:"date"`
	AccessKey string                 `json:"access_key"`
	PostID    int                    `json:"post_id"`
}

type AttachmentVideo struct {
	AccessKey     string                 `json:"access_key"`
	CanComment    int                    `json:"can_comment"`
	CanEdit       int                    `json:"can_edit"`
	CanLike       int                    `json:"can_like"`
	CanRepost     int                    `json:"can_repost"`
	CanSubscribe  int                    `json:"can_subscribe"`
	CanAddToFaves int                    `json:"can_add_to_faves"`
	CanAdd        int                    `json:"can_add"`
	CanAttachLink int                    `json:"can_attach_link"`
	Comments      int                    `json:"comments"`
	Date          int64                  `json:"date"`
	Description   string                 `json:"description"`
	Duration      int                    `json:"duration"`
	Image         []*AttachmentImageInfo `json:"image"`
	FirstFrame    []*AttachmentImageInfo `json:"first_frame"`
	Width         int                    `json:"width"`
	Height        int                    `json:"height"`
	ID            int                    `json:"id"`
	OwnerID       int                    `json:"owner_id"`
	UserID        int                    `json:"user_id"`
	Title         string                 `json:"title"`
	IsFavorite    int                    `json:"is_favorite"`
	TrackCode     string                 `json:"track_code"`
	Type          string                 `json:"type"`
	Views         int                    `json:"views"`
	LocalViews    int                    `json:"local_views"`
	Platform      string                 `json:"platform"`
}

type AttachmentAudio struct {
	Artist     string                   `json:"artist"`
	ID         int                      `json:"id"`
	OwnerID    int                      `json:"owner_id"`
	Title      string                   `json:"title"`
	Duration   int                      `json:"duration"`
	AccessKey  string                   `json:"access_key"`
	IsLicensed bool                     `json:"is_licensed"`
	TrackCode  string                   `json:"track_code"`
	Url        string                   `json:"url"`
	Date       int64                    `json:"date"`
	Album      *AttachmentAudioAlbum    `json:"album"`
	MainArtist []*AttachmentAudioArtist `json:"main_artist"`
}

type AttachmentAudioAlbum struct {
	ID        int                        `json:"id"`
	Title     string                     `json:"title"`
	OwnerID   int                        `json:"owner_id"`
	AccessKey string                     `json:"access_key"`
	Thumb     *AttachmentAudioAlbumThumb `json:"thumb"`
}

type AttachmentAudioAlbumThumb struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Photo34  string `json:"photo_34"`
	Photo68  string `json:"photo_68"`
	Photo135 string `json:"photo_135"`
	Photo270 string `json:"photo_270"`
	Photo300 string `json:"photo_300"`
	Photo600 string `json:"photo_600"`
}

type AttachmentAudioArtist struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	ID     string `json:"id"`
}

type AttachmentDoc struct {
	ID         int    `json:"id"`
	OwnerID    int    `json:"owner_id"`
	Title      string `json:"title"`
	Size       int    `json:"size"`
	Ext        string `json:"ext"`
	Url        string `json:"url"`
	Date       int64  `json:"date"`
	Type       int    `json:"type"`
	IsLicensed int    `json:"is_licensed"` // NOTE(Pedro): There is some inconsistency on VK API, this field can return "true" or zero for false
	AccessKey  string `json:"access_key"`
}

type AttachmentSticker struct {
	ProductID            int                    `json:"product_id"`
	StickerID            int                    `json:"sticker_id"`
	Images               []*AttachmentImageInfo `json:"images"`
	ImagesWithBackground []*AttachmentImageInfo `json:"images_with_background"`
}

type Poll struct {
	ID        int           `json:"id"`
	OwnerID   int           `json:"owner_id"`
	Created   int64         `json:"created"`
	Question  string        `json:"question"`
	Votes     int           `json:"votes"`
	Answers   []*PollAnswer `json:"answers"`
	Anonymous bool          `json:"anonymous"`
	Multiple  bool          `json:"multiple"`
	AnswerIDS []int         `json:"answer_ids"`
	EndDate   int64         `json:"end_date"`
	Closed    bool          `json:"closed"`
	IsBoard   bool          `json:"is_board"`
	CanEdit   bool          `json:"can_edit"`
	CanVote   bool          `json:"can_vote"`
	CanReport bool          `json:"can_report"`
	CanShare  bool          `json:"can_share"`
}

type PollAnswer struct {
	ID    int     `json:"id"`
	Text  string  `json:"text"`
	Votes int     `json:"votes"`
	Rate  float64 `json:"rate"`
}

type Comments struct {
	Count      int             `json:"count"`
	Comments   []*TopicComment `json:"items"`
	Poll       *Poll           `json:"poll"`
	Profiles   []*User         `json:"profiles"`
	RealOffset int             `json:"real_offset"`
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
	err = json.Unmarshal(resp.Response, &comments)
	if err != nil {
		return nil, err
	}

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
	err = json.Unmarshal(resp.Response, &topics)
	if err != nil {
		return nil, err
	}

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
