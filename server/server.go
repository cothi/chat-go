package server

import (
	"fmt"
	"net"
	"time"

	"github.com/cothi/chat-go/utils"
)

type Message struct {
	time time.Time
	text string
}

type room struct {
	room_name string
	port      []string
	clients   []*Client
	message   []string
}

var (
	room_1 room
)

func (r *room) broadcast(msg []byte) {
	for _, c := range r.clients {
		c.box <- msg
	}
}

type Client struct {
	name string
	conn net.Conn
	box  chan []byte
}

func InitClient(conn net.Conn) {
	client := Client{
		name: "anonymous",
		conn: conn,
		box:  make(chan []byte),
	}
	go client.Write()
	go client.Read()
	room_1.clients = append(room_1.clients, &client)
}

func (c *Client) Write() {
	for {

		s := <-c.box
		c.conn.Write(s)
		fmt.Println(c.name, c.conn.LocalAddr(), string(s))
	}
}

func (c *Client) Read() {
	recv := make([]byte, 4096)
	for {
		n, e := c.conn.Read(recv)
		utils.Error_check(e)
		room_1.broadcast(recv[:n])
	}
}

func Start(port string) {

	l, err := net.Listen("tcp", port)
	utils.Error_check(err)
	defer l.Close()

	fmt.Printf("Start Server port %s \n", port)
	for {
		conn, err := l.Accept()
		InitClient(conn)
		utils.Error_check(err)
	}
}
