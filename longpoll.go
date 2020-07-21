package vkapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type LongPollServer struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	TS     int64  `json:"ts"`
}

type LongPollUpdate struct {
	Failed  int             `json:"failed"`
	TS      int64           `json:"ts"`
	Updates [][]interface{} `json:"updates"`
}

type BotsLongPollServer struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	TS     string `json:"ts"`
}

type BotsLongPollUpdate struct {
	Failed  int                 `json:"failed"`
	TS      string              `json:"ts"`
	Updates []BotsLongPollEvent `json:"updates"`
}

type LongPollMessage struct {
	MessageType  string
	MessageID    int
	MessageFlags int
	UserID       int64
	Date         int64
	Title        string
	Body         string
	Attachments  map[string]string
}

type BotsLongPollEvent struct {
	Type   string             `json:"type"`
	Object BotsLongPollObject `json:"object"`
}

type BotsLongPollObject struct {
	Message    BotlLongPollDM `json:"message"`
	ClientInfo interface{}    `json:"client_info`
}

type BotlLongPollDM struct {
	MessageID    int                 `json:"id"`
	Date         int64               `json:"date"`
	PeerID       int64               `json:"peer_id"`
	SendByID     int64               `json:"from_id"`
	Text         string              `json:"text"`
	RandomID     int64               `json:"random_id"`
	Attachments  []MessageAttachment `json:"attachments"`
	Isimportant  bool                `json:"important"`
	FwdMessages  []BotlLongPollDM    `json:"fwd_messages"`
	ReplyMessage *BotlLongPollDM     `json:"reply_message"`
	Payload      string              `json:"payload"`
}

func (client *VKClient) getLongPollServer() (LongPollServer, error) {
	resp, err := client.MakeRequest("messages.getLongPollServer", nil)
	if err != nil {
		return LongPollServer{}, err
	}
	var server LongPollServer
	err = json.Unmarshal(resp.Response, &server)
	if err != nil {
		return LongPollServer{}, err
	}

	return server, nil
}

func (client *VKClient) longpollRequest(server LongPollServer) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://"+server.Server, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("act", "a_check")
	q.Add("key", server.Key)
	q.Add("ts", strconv.FormatInt(server.TS, 10))
	q.Add("wait", "25")
	q.Add("mode", "2")
	q.Add("version", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *VKGroupBot) getBotsLongPollServer() (BotsLongPollServer, error) {
	params := url.Values{}
	params.Set("group_id", strconv.FormatInt(int64(client.ID), 10))
	resp, err := client.MakeRequest("groups.getLongPollServer", params)
	if err != nil {
		return BotsLongPollServer{}, err
	}
	var server BotsLongPollServer
	err = json.Unmarshal(resp.Response, &server)
	if err != nil {
		return BotsLongPollServer{}, err
	}

	return server, nil
}

func (client *VKGroupBot) botsLongpollRequest(server BotsLongPollServer) ([]byte, error) {
	req, err := http.NewRequest("GET", server.Server, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("act", "a_check")
	q.Add("key", server.Key)
	q.Add("ts", server.TS)
	q.Add("wait", "25")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *VKClient) ListenLongPollServerWithCancel(cancelCtx context.Context) {
	server, err := client.getLongPollServer()
	if err != nil {
		fmt.Println("failed to get longpoll server")
		return
	}

	for {
		select {
		case <-cancelCtx.Done():
			return
		default:
		}

		body, err := client.longpollRequest(server)
		if err != nil {
			log.Printf("longpoll request failed: %s", err)
			time.Sleep(time.Second * 5)
			continue
		}

		var updates LongPollUpdate
		err = json.Unmarshal(body, &updates)

		switch updates.Failed {
		case 0:
			for _, update := range updates.Updates {
				updateID := update[0].(float64)
				message := new(LongPollMessage)

				switch updateID {
				case 4: //new message
					message.MessageID = int(update[1].(float64))
					message.MessageFlags = int(update[2].(float64))
					message.UserID = int64(update[3].(float64))
					message.Date = int64(update[4].(float64))
					message.Title = update[5].(string)
					message.Body = update[6].(string)
					message.Attachments = make(map[string]string)

					for k, v := range update[7].(map[string]interface{}) {
						message.Attachments[k] = v.(string)
					}

					if message.MessageFlags == 19 || message.MessageFlags == 51 ||
						message.MessageFlags == 531 || message.MessageFlags == 563 ||
						message.MessageFlags == 3 || message.MessageFlags == 35 {
						message.MessageType = "msgout"
					} else {
						message.MessageType = "msgin"
					}

				case 2: //message deleted
					message.MessageType = "msgdel"
					message.MessageID = int(update[1].(float64))
					message.UserID = int64(update[3].(float64))
				case 3: //message read
					message.MessageType = "msgread"
					message.MessageID = int(update[1].(float64))
				case 8: //user online
					message.MessageType = "msgonline"
					message.UserID = int64(update[1].(float64))
				}
				client.handleCallback(message.MessageType, message)
			}
			server.TS = updates.TS

		case 1:
			server.TS = updates.TS

		case 2:
			newSrv, err := client.getLongPollServer()
			if err != nil {
				log.Printf("failed to get longpoll server for key update: %s", err)
				continue
			}
			server.Key = newSrv.Key

		case 3:
			newSrv, err := client.getLongPollServer()
			if err != nil {
				log.Printf("failed to get longpoll for key/ts update: %s", err)
				continue
			}
			server.Key = newSrv.Key
			server.TS = newSrv.TS
		}
	}
}

func (client *VKGroupBot) ListenBotsLongPollServerWithCancel(cancelCtx context.Context) {
	server, err := client.getBotsLongPollServer()
	if err != nil {
		fmt.Println("failed to get botsLongpoll server")
		return
	}

	for {
		select {
		case <-cancelCtx.Done():
			return
		default:
		}
		body, err := client.botsLongpollRequest(server)
		if err != nil {
			log.Printf("botsLongpoll request failed: %s", err)
			time.Sleep(time.Second * 5)
			continue
		}

		var updates BotsLongPollUpdate
		err = json.Unmarshal(body, &updates)
		switch updates.Failed {
		case 0:
			for _, update := range updates.Updates {
				switch update.Type {
				case "message_new", "message_reply", "message_edit":
					client.handleBotsCallback(update.Type, &update.Object)
				default:
					log.Printf("botsLongpoll update with non 'message' type: %s", update.Type)
				}
			}
			server.TS = updates.TS
		case 1:
			server.TS = updates.TS

		case 2:
			newSrv, err := client.getBotsLongPollServer()
			if err != nil {
				log.Printf("failed to get botsLongpoll server for key update: %s", err)
				continue
			}
			server.Key = newSrv.Key

		case 3:
			newSrv, err := client.getBotsLongPollServer()
			if err != nil {
				log.Printf("failed to get botsLongpoll for key/ts update: %s", err)
				continue
			}
			server.Key = newSrv.Key
			server.TS = newSrv.TS
		}
	}
}

func (client *VKClient) ListenLongPollServer() {
	client.ListenLongPollServerWithCancel(context.Background())
}

func (client *VKGroupBot) ListenBotsLongPollServer() {
	client.ListenBotsLongPollServerWithCancel(context.Background())
}
