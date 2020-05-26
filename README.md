# VK API wrapper in golang

## Installing

```
go get github.com/himidori/golang-vk-api
```

## Authorizing using username and password

```go
client, err := vkapi.NewVKClient(vkapi.DeviceIPhone, "username", "password")
```

## Authorizing using access token

```go
client, err := vkapi.NewVKClientWithToken("token", nil)
```

## Listening longpoll events

```go
// listening received messages
client.AddLongpollCallback("msgin", func(m *vkapi.LongPollMessage) {
	fmt.Printf("new message received from uid %d\n", m.UserID)
})

// listening deleted messages
client.AddLongpollCallback("msgdel", func(m *vkapi.LongPollMessage) {
	fmt.Printf("message %d was deleted\n", m.MessageID)
})

// listening sent messages
client.AddLongpollCallback("msgout", func(m *vkapi.LongPollMessage) {
	fmt.Printf("sent message to uid %d\n", m.UserID)
})

// listening read messages
client.AddLongpollCallback("msgread", func(m *vkapi.LongPollMessage) {
	fmt.Printf("message %d was read\n", m.MessageID)
})

// listening users online
client.AddLongpollCallback("msgonline", func(m *vkapi.LongPollMessage) {
	fmt.Printf("user %d is now online\n", m.UserID)
})

// starting 
client.ListenLongPollServer()
```
