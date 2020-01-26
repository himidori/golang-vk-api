package vkapi

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"
)

// NOTE(Pedro): Fill the fields below to be able to run the integration tests
// it will create/delete/edit/fix/unfix topcis and comments, use on a test community
const (
	testUser      = "" // email
	testPassword  = "" // password
	testCommunity = 0  // community id
)

var client *VKClient

func TestMain(m *testing.M) {
	var err error
	client, err = NewVKClient(DeviceIPhone, testUser, testPassword)
	if err != nil {
		fmt.Println(err)
		return
	}

	m.Run()
}

func TestVKClient_BoardAddTopic(t *testing.T) {
	id, err := client.BoardAddTopic(testCommunity, "This is the title", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if id == 0 {
		t.Fatalf("got wrong id %d", id)
	}
}

func TestVKClient_BoardCloseTopic(t *testing.T) {
	id, err := client.BoardAddTopic(testCommunity, "This is the title of a closed topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if id == 0 {
		t.Fatalf("got wrong id %d", id)
	}

	closed, err := client.BoardCloseTopic(testCommunity, id)
	if err != nil {
		t.Fatal(err)
	}

	if !closed {
		t.Fatal("could not close the topic")
	}
}

func TestVKClient_BoardCreateComment(t *testing.T) {
	id, err := client.BoardAddTopic(testCommunity, "This is the title of a topic with comments", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if id == 0 {
		t.Fatalf("got wrong id %d", id)
	}

	commentID, err := client.BoardCreateComment(testCommunity, id, "test comment", []string{}, false, 0)
	if err != nil {
		t.Fatal(err)
	}

	if commentID == 0 {
		t.Fatalf("got wrong comment id %d", commentID)
	}

	commentIDWithSticker, err := client.BoardCreateComment(testCommunity, id, "", []string{}, false, 2466)
	if err != nil {
		t.Fatal(err)
	}

	if commentIDWithSticker == 0 {
		t.Fatalf("got wrong sticker comment id %d", commentIDWithSticker)
	}
}

func TestVKClient_BoardDeleteComment(t *testing.T) {
	topicID, err := client.BoardAddTopic(testCommunity, "This is the title of a deleted topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if topicID == 0 {
		t.Fatalf("got wrong id %d", topicID)
	}

	commentID, err := client.BoardCreateComment(testCommunity, topicID, "test comment", []string{}, false, 0)
	if err != nil {
		t.Fatal(err)
	}

	if commentID == 0 {
		t.Fatalf("got wrong comment id %d", commentID)
	}

	deleted, err := client.BoardDeleteComment(testCommunity, topicID, commentID)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Fatal("could not delete the comment")
	}
}

func TestVKClient_BoardDeleteTopic(t *testing.T) {
	topicID, err := client.BoardAddTopic(testCommunity, "This is the title of a deleted topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if topicID == 0 {
		t.Fatalf("got wrong id %d", topicID)
	}

	deleted, err := client.BoardDeleteTopic(testCommunity, topicID)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Fatal("could not delete the topic")
	}
}

func TestVKClient_BoardEditComment(t *testing.T) {
	topicID, err := client.BoardAddTopic(testCommunity, "This is the title of a test topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if topicID == 0 {
		t.Fatalf("got wrong id %d", topicID)
	}

	commentID, err := client.BoardCreateComment(testCommunity, topicID, "test comment", []string{}, false, 0)
	if err != nil {
		t.Fatal(err)
	}

	if commentID == 0 {
		t.Fatalf("got wrong comment id %d", commentID)
	}

	edited, err := client.BoardEditComment(testCommunity, topicID, commentID, "changed", []string{})
	if err != nil {
		t.Fatal(err)
	}

	if !edited {
		t.Fatal("could not edit the comment")
	}

	comments, err := client.BoardGetComments(testCommunity, topicID, 2, url.Values{})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments.Comments) != 2 {
		t.Fatal("got wrong number of comments")
	}

	// NOTE(Pedro): we test the second comment, since the first one is the topic message
	if comments.Comments[1].Text != "changed" {
		t.Fatalf("the message was not updated, expected %s, got %s", "changed", comments.Comments[0].Text)
	}
}

func TestVKClient_BoardEditTopic(t *testing.T) {
	topicID, err := client.BoardAddTopic(testCommunity, "This is the title of a test topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if topicID == 0 {
		t.Fatalf("got wrong id %d", topicID)
	}

	edited, err := client.BoardEditTopic(testCommunity, topicID, "New title")
	if err != nil {
		t.Fatal(err)
	}

	if !edited {
		t.Fatal("could not edit the topic")
	}

	var params = url.Values{}
	params.Add("topic_ids", strconv.Itoa(topicID))

	topics, err := client.BoardGetTopics(testCommunity, 1, params)
	if err != nil {
		t.Fatal(err)
	}

	if len(topics.Topics) != 1 {
		t.Fatal("got wrong number of topics")
	}

	if topics.Topics[0].Title != "New title" {
		t.Fatalf("got wrong title, expected %s, got %s", "New title", topics.Topics[0].Title)
	}
}

func TestVKClient_BoardFixTopic(t *testing.T) {
	topicID, err := client.BoardAddTopic(testCommunity, "This is the title of a fixed topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if topicID == 0 {
		t.Fatalf("got wrong id %d", topicID)
	}

	fixed, err := client.BoardFixTopic(testCommunity, topicID)
	if err != nil {
		t.Fatal(err)
	}

	if !fixed {
		t.Fatal("could not fix the topic")
	}
}

func TestVKClient_BoardGetComments(t *testing.T) {
	topicID, err := client.BoardAddTopic(testCommunity, "This is the title of a test topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if topicID == 0 {
		t.Fatalf("got wrong id %d", topicID)
	}

	commentID, err := client.BoardCreateComment(testCommunity, topicID, "test comment", []string{}, false, 0)
	if err != nil {
		t.Fatal(err)
	}

	if commentID == 0 {
		t.Fatalf("got wrong comment id %d", commentID)
	}

	comments, err := client.BoardGetComments(testCommunity, topicID, 100, url.Values{})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments.Comments) != 2 {
		t.Fatal("got wrong number of comments")
	}
}

func TestVKClient_BoardGetTopics(t *testing.T) {
	topics, err := client.BoardGetTopics(testCommunity, 55, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(topics.Topics) != 55 {
		t.Fatal("got wrong number of topics")
	}

	fmt.Println(topics.Topics[0].Title)

	comments, err := client.BoardGetComments(testCommunity, topics.Topics[0].ID, 100, nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(comments.Comments[0].Text)
}

func TestVKClient_BoardOpenTopic(t *testing.T) {
	id, err := client.BoardAddTopic(testCommunity, "This is the title of a open topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if id == 0 {
		t.Fatalf("got wrong id %d", id)
	}

	closed, err := client.BoardCloseTopic(testCommunity, id)
	if err != nil {
		t.Fatal(err)
	}

	if !closed {
		t.Fatal("could not close the topic")
	}

	opened, err := client.BoardOpenTopic(testCommunity, id)
	if err != nil {
		t.Fatal(err)
	}

	if !opened {
		t.Fatal("could not open the topic")
	}
}

func TestVKClient_BoardRestoreComment(t *testing.T) {
	topicID, err := client.BoardAddTopic(testCommunity, "This is the title of a restored topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if topicID == 0 {
		t.Fatalf("got wrong id %d", topicID)
	}

	commentID, err := client.BoardCreateComment(testCommunity, topicID, "test comment to be restored", []string{}, false, 0)
	if err != nil {
		t.Fatal(err)
	}

	if commentID == 0 {
		t.Fatalf("got wrong comment id %d", commentID)
	}

	deleted, err := client.BoardDeleteComment(testCommunity, topicID, commentID)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Fatal("could not delete the comment")
	}

	restored, err := client.BoardRestoreComment(testCommunity, topicID, commentID)
	if err != nil {
		t.Fatal(err)
	}

	if !restored {
		t.Fatal("could not restore the comment")
	}
}

func TestVKClient_BoardUnfixTopic(t *testing.T) {
	topicID, err := client.BoardAddTopic(testCommunity, "This is the title of a unfixed topic", "This is the text", false, []string{})
	if err != nil {
		t.Fatal(err)
	}

	if topicID == 0 {
		t.Fatalf("got wrong id %d", topicID)
	}

	fixed, err := client.BoardFixTopic(testCommunity, topicID)
	if err != nil {
		t.Fatal(err)
	}

	if !fixed {
		t.Fatal("could not fix the topic")
	}

	unfixed, err := client.BoardUnfixTopic(testCommunity, topicID)
	if err != nil {
		t.Fatal(err)
	}

	if !unfixed {
		t.Fatal("could not unfix the topic")
	}
}
