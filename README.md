# VK API wrapper in golang

## Installing
go get github.com/himidori/golang-vk-api

## Example of usage

```go
package main

import (
	"log"
	"net/url"

	"github.com/himidori/golang-vk-api"
)

func main() {
	//authorizing as IPhone device using login and password
	client, err := vkapi.NewVKClient(vkapi.DeviceIPhone, "", "")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("authorized as %s %s (id%d)\n", client.Self.FirstName, client.Self.LastName,
		client.Self.UID)

	//sending a post with photo attachments to a group's wall
	pics := []string{
		"/home/yuimaestro/agirl.jpg",
		"/home/yuimaestro/agirl2.jpg",
	}

	photos, err := client.UploadGroupWallPhotos(111, pics)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("uploaded %d pictures\n", len(photos))

	s := client.GetPhotosString(photos)
	params := url.Values{}
	params.Set("attachments", s)

	msg, err := client.WallPost(-111, "testmessage", params)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("successfuly posted a message with ID %d\n", msg)

	//starting a longpoll channel to listen for events
	ch, err := client.ListenLongPollServer()
	if err != nil {
		log.Fatal(err)
	}

	for msg := range ch {
		log.Printf("new message from user %d: %s\n", msg.UserID, msg.Body)
	}

}
