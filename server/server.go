package server

import (
	"fmt"
	"net"
	"time"

	"github.com/cothi/chat-go/utils"
)

var (
	Lobby Chan
)

type Message struct {
	time time.Time
	text string
}

type Chan struct {
	room map[string]*room
}

func (c *Chan) Init() {
	c.room["lobby"] = &room{
		room_name: "lobby",
	}
}

func (c *Chan) Create(name string) {
	c.room[name] = &room{}
}

type room struct {
	room_name string
	port      []string
	clients   []*Client
	message   []string
}

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
	Lobby.room["lobby"].clients = append(Lobby.room["lobby"].clients, &client)
}

func (c *Client) Write() {
	for {
		s, ok := <-c.box
		if !ok {
			fmt.Printf("Close client %s %s \n", c.name, c.conn.LocalAddr())
			return
		}
		c.conn.Write(s)
		fmt.Println(c.name, c.conn.LocalAddr(), string(s))
	}
}

func (c *Client) Read() {
	recv := make([]byte, 4096)
	for {
		n, e := c.conn.Read(recv)

		if e != nil {
			c.conn.Close()
			close(c.box)
			return
		}
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
