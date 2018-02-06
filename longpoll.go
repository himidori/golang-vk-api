package vkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	MessageType  int
	MessageID    int
	MessageFlags int
	UserID       int
	Date         int64
	Title        string
	Body         string
	Attachments  map[string]string
}

type LongPollChannel <-chan LongPollMessage

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

func (client *VKClient) longpollGet(server LongPollServer) ([]byte, error) {
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

func (client *VKClient) retry(server LongPollServer) ([]byte, error) {
	var resp []byte
	var err error

	for i := 0; i < 3; i++ {
		resp, err = client.longpollGet(server)

		if err == nil {
			return resp, nil
		}

		fmt.Printf("request to longpoll server failed, " +
			"sleeping for 3 secs\n")
		time.Sleep(time.Second * 3)
	}

	return nil, err
}

func (client *VKClient) ListenLongPollServer() (LongPollChannel, error) {
	ch := make(chan LongPollMessage, 10)
	server, err := client.getLongPollServer()
	if err != nil {
		return ch, err
	}

	go func() {
		for {
			body, err := client.retry(server)
			if err != nil {
				fmt.Println("failed request to longpoll server after 3 retries")
				close(ch)
				return
			}

			var updates LongPollUpdate
			err = json.Unmarshal(body, &updates)

			switch updates.Failed {
			case 0:
				for _, update := range updates.Updates {
					updateID := update[0].(float64)

					switch updateID {
					case 4: //new message
						var message LongPollMessage
						message.MessageType = 4
						message.MessageID = int(update[1].(float64))
						message.MessageFlags = int(update[2].(float64))
						message.UserID = int(update[3].(float64))
						message.Date = int64(update[4].(float64))
						message.Title = update[5].(string)
						message.Body = update[6].(string)
						message.Attachments = make(map[string]string)

						for k, v := range update[7].(map[string]interface{}) {
							message.Attachments[k] = v.(string)
						}

						ch <- message

					case 2: //message deleted
						var message LongPollMessage
						message.MessageType = 2
						message.MessageID = int(update[1].(float64))
						message.UserID = int(update[3].(float64))
						ch <- message
					}
				}
				server.TS = updates.TS
			case 1:
				server.TS = updates.TS
			case 2, 3:
				server, err = client.getLongPollServer()
				if err != nil {
					fmt.Println("error requesting longpoll server")
				}
			}
		}
	}()

	return ch, nil
}
