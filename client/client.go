package client

import (
	"encoding/json"
	"net"
	"strings"
	"sync"

	"github.com/cothi/chat-go/utils"
)

type MessageKind int

// Post structure
type Post struct {
	Chat     string      `json:"chat"`
	Time     string      `json:"time"`
	Roomname string      `json:"roomname"`
	Nickname string      `json:"nickname"`
	Kind     MessageKind `json:"kind"`
}

// Client structure
type Client struct {
	Name   string
	Conn   net.Conn
	Outbox chan []byte
	Inbox  chan []byte
	m      sync.Mutex
}

const (
	MessageCreateMsg MessageKind = iota
	MessageCreateRoom
	MessageJoinRoom
)

// read
func (c *Client) Read() {
	recv := make([]byte, 4096)
	for {
		i, e := c.Conn.Read(recv)
		utils.Error_check(e)
		c.Inbox <- recv[:i]
	}
}

// write
func (c *Client) Write() {
	for {
		m, ok := <-c.Outbox
		if !ok {
			return
		}
		c.Conn.Write(m)
	}
}

// client
func (c *Client) ClientInit(serverPort string) *Client {
	conn, _ := net.Dial("tcp", ":"+serverPort)
	c.Name = "anonymous"
	c.Conn = conn
	c.Outbox = make(chan []byte)
	c.Inbox = make(chan []byte)
	go c.Read()
	go c.Write()
	return c
}

// send handle message
func (c *Client) SendHandleMessage(p *Post) {
	splitWord := strings.Split(p.Chat, " ")

	switch splitWord[0] {
	case "/join":
		p.Kind = MessageJoinRoom
		p.Chat = splitWord[1]
	case "/create":
		p.Kind = MessageCreateRoom
		p.Chat = splitWord[1]
	default:
		p.Kind = MessageCreateMsg
	}
	b, e := json.Marshal(p)
	utils.Error_check(e)
	c.Outbox <- b
}
