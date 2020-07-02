package main

import (
	"net/url"

	vkapi "github.com/himidori/golang-vk-api"
)

/*
 * В данном примере рассматривается публикация
 * поста с картинкой на стену указанного сообщества
 * используя ссылку на существующее изображение
 * Примечание:
 * На текущий момент функция поддерживает лишь единоразовую загрузку файла
 */

func publicateToWallbyURL() {
	var vkToken string
	var vkGroupID int
	var urlImage string
	client, err := vkapi.NewVKClientWithToken(vkToken, &vkapi.TokenOptions{})
	if err != nil {
		panic(err)
	}
	photo, err := client.UploadByLinkGroupWallPhotos(vkGroupID, urlImage)
	if err != nil {
		panic(err)
	}
	params := url.Values{}
	params.Set("attachments", client.GetPhotosString(photo))
	client.WallPost(vkGroupID, "message", params)
}

/*
 * В данном примере рассматривается публикация
 * поста с картинкой на стену указанного сообщества
 * используя существующее изображение на локальном диске
 * Примечание:
 * Согласно API возможно загрузить только шесть изображений, общим весом до 50 мб
 */
func publicateToWallbyFile() {
	var vkToken string
	var vkGroupID int

	client, err := vkapi.NewVKClientWithToken(vkToken, &vkapi.TokenOptions{})
	if err != nil {
		panic(err)
	}

	files := []string{"example.png"}

	photo, err := client.UploadGroupWallPhotos(vkGroupID, files)
	if err != nil {
		panic(err)
	}
	params := url.Values{}
	params.Set("attachments", client.GetPhotosString(photo))
	client.WallPost(vkGroupID, "message", params)
}
