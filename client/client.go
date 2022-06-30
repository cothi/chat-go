package client

import (
	"net"
	"sync"

	"github.com/cothi/chat-go/utils"
)

type Client struct {
	Name   string
	Conn   net.Conn
	Outbox chan []byte
	Inbox  chan []byte
	m      sync.Mutex
}

func (c *Client) Read() {
	recv := make([]byte, 4096)
	for {
		i, e := c.Conn.Read(recv)
		utils.Error_check(e)
		c.Inbox <- recv[:i]
	}
}

func (c *Client) Write() {
	for {
		m, ok := <-c.Outbox
		if !ok {
			return
		}
		c.Conn.Write(m)
	}
}

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
