package client

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/cothi/tcp-chat-remodel/utils"
)

type Client struct {
	Conn   net.Conn
	Outbox chan []byte
	Inbox  chan string
	m      sync.Mutex
}

func (c *Client) Read() {

	recv := make([]byte, 4096)

	for {
		n, e := c.Conn.Read(recv)
		// fmt.Println(n)
		utils.Error_check(e)

		str := fmt.Sprintf("%s: %s", c.Conn.LocalAddr(), strings.TrimSpace(string(recv[:n])))
		c.Inbox <- str
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

func ClientInit() *Client {
	conn, err := net.Dial("tcp", ":8000")
	utils.Error_check(err)

	c := &Client{
		Conn:   conn,
		Outbox: make(chan []byte),
		Inbox:  make(chan string),
	}

	go c.Read()
	go c.Write()

	return c
}

// Reads from Stdin, and outputs to the socket.
