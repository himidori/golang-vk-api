package vkapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func (client *VKClient) getLongPollServer() (LongPollServer, error) {
	resp, err := client.makeRequest("messages.getLongPollServer", nil)
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
			log.Println("longpoll request failed: %s", err)
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
				log.Println("failed to get longpoll server for key update: %s", err)
				continue
			}
			server.Key = newSrv.Key

		case 3:
			newSrv, err := client.getLongPollServer()
			if err != nil {
				log.Println("failed to get longpoll for key/ts update: %s", err)
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