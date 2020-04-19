package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type VideoImage struct {
	Url         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	WithPadding int    `json:"with_padding"`
}

type VideoFirstFrame struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Video struct {
	ID            int                `json:"id"`
	OwnerID       int                `json:"owner_id"`
	UserID        int                `json:"user_id"`
	Title         string             `json:"title"`
	Duration      int                `json:"duration"`
	Description   string             `json:"description"`
	Date          int64              `json:"date"`
	Comments      int                `json:"comments"`
	Views         int                `json:"views"`
	Photo130      string             `json:"photo_130"`
	Photo320      string             `json:"photo_320"`
	Photo640      string             `json:"photo_640"`
	IsFavorite    bool               `json:"is_favorite"`
	AddingDate    int64              `json:"adding_date"`
	Player        string             `json:"player"`
	Platftorm     string             `json:"platform"`
	CanAdd        int                `json:"can_add"`
	CanComment    int                `json:"can_comment"`
	CanLike       int                `json:"can_like"`
	CanRepost     int                `json:"can_repost"`
	CanSubscribe  int                `json:"can_subscribe"`
	CanAddToFaves int                `json:"can_add_to_faves"`
	Image         []*VideoImage      `json:"image"`
	FirstFrame    []*VideoFirstFrame `json:"first_frame"`
	Type          string             `json:"type"`
	Likes         struct {
		UserLikes int `json:"user_likes"`
		Count     int `json:"count"`
	} `json:"likes"`
	Reposts struct {
		Count        int `json:"count"`
		UserReposted int `json:"user_reposted"`
	} `json:"reposts"`
	Files struct {
		External string `json:"external"`
		FLV_320  string `json:"mp4_320"`
		MP4_240  string `json:"mp4_240"`
		MP4_360  string `json:"mp4_360"`
		MP4_480  string `json:"mp4_480"`
		MP4_720  string `json:"mp4_720"`
		MP4_1080 string `json:"mp4_1080"`
	} `json:"files"`
}

type Videos struct {
	Count    int      `json:"count"`
	Videos   []*Video `json:"items"`
	Profiles []*User  `json:"profiles"`
	Groups   []*Group `json:"groups"`
}

type VideoAttachment struct {
	ID            int           `json:"id"`
	OwnerID       int           `json:"owner_id"`
	Title         string        `json:"title"`
	Duration      int           `json:"duration"`
	Description   string        `json:"description"`
	Date          int64         `json:"date"`
	AddingDate    int64         `json:"adding_date"`
	Views         int           `json:"views"`
	Width         int           `json:"width"`
	Height        int           `json:"height"`
	Photo130      string        `json:"photo130"`
	Photo320      string        `json:"photo320"`
	Photo800      string        `json:"photo800"`
	FirstFrame320 string        `json:"first_frame_320"`
	FirstFrame160 string        `json:"first_frame_160"`
	FirstFrame130 string        `json:"first_frame_130"`
	FirstFrame800 string        `json:"first_frame_800"`
	Player        string        `json:"player"`
	CanEdit       int           `json:"can_edit"`
	CanAdd        int           `json:"can_add"`
	CanComment    int           `json:"can_comment"`
	CanLike       int           `json:"can_like"`
	CanRepost     int           `json:"can_repost"`
	CanSubscribe  int           `json:"can_subscribe"`
	CanAddToFaves int           `json:"can_add_to_faves"`
	AccessKey     string        `json:"access_key"`
	Image         []*VideoImage `json:"image"`
}

func (client *VKClient) VideoGet(ownerID int, count int, params url.Values) (*Videos, error) {
	if params == nil {
		params = url.Values{}
	}

	params.Set("owner_id", strconv.Itoa(ownerID))
	params.Set("count", strconv.Itoa(count))

	resp, err := client.MakeRequest("video.get", params)
	if err != nil {
		return nil, err
	}

	var videos *Videos

	if err = json.Unmarshal(resp.Response, &videos); err != nil {
		return nil, err
	}

	return videos, nil
}
