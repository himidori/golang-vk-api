package vkapi

import (
	"sync"
)

type callbackHandler struct {
	sync.Mutex
	events map[string]func(*LongPollMessage)
}

type botsCallBackHandler struct {
	sync.Mutex
	events map[string]func(*BotsLongPollObject)
}

func (client *VKClient) AddLongpollCallback(name string, f func(*LongPollMessage)) {
	client.cb.Lock()
	if _, exists := client.cb.events[name]; !exists {
		client.cb.events[name] = f
	}
	client.cb.Unlock()
}

func (client *VKClient) DeleteLongpollCallback(name string) {
	client.cb.Lock()
	if _, exists := client.cb.events[name]; exists {
		delete(client.cb.events, name)
	}
	client.cb.Unlock()
}

func (client *VKClient) handleCallback(name string, msg *LongPollMessage) {
	client.cb.Lock()
	_, exists := client.cb.events[name]
	client.cb.Unlock()

	if exists {
		client.cb.events[name](msg)
	}
}

func (client *VKGroupBot) AddBotsLongpollCallback(name string, f func(*BotsLongPollObject)) {
	client.cb.Lock()
	if _, exists := client.cb.events[name]; !exists {
		client.cb.events[name] = f
	}
	client.cb.Unlock()
}

func (client *VKGroupBot) DeleteBotsLongpollCallback(name string) {
	client.cb.Lock()
	if _, exists := client.cb.events[name]; exists {
		delete(client.cb.events, name)
	}
	client.cb.Unlock()
}

func (client *VKGroupBot) handleBotsCallback(name string, msg *BotsLongPollObject) {
	client.cb.Lock()
	_, exists := client.cb.events[name]
	client.cb.Unlock()

	if exists {
		client.cb.events[name](msg)
	}
}
