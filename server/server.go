package server

import (
	"encoding/json"
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

	c.room = make(map[string]*room)
	room1 := room{
		room_name: "0",
	}

	c.room["0"] = &room1

}

func (c *Chan) Create(name string) {
	c.room[name] = &room{}
}

type room struct {
	room_name string
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
	Lobby.room["0"].clients = append(Lobby.room["0"].clients, &client)
	go client.Write()
	go client.Read()
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
	var msg Post
	for {
		i, e := c.conn.Read(recv)
		fmt.Println(recv[:i])

		if e != nil {
			c.conn.Close()
			close(c.box)
			return
		}
		json.Unmarshal(recv[:i], &msg)
		fmt.Println(msg.RoomNum)
    Lobby.room[msg.RoomNum].broadcast(recv[:i])
	}
}

func Start(port string) {

	Lobby.Init()
	l, err := net.Listen("tcp", port)
	utils.Error_check(err)
	defer l.Close()
	fmt.Printf("Start Server port %s \n", port)
	for {
		conn, err := l.Accept()
		utils.Error_check(err)
		InitClient(conn)
	}
}

type Post struct {
	Chat    string `json:"chat"`
	Time    string `json:"time"`
	RoomNum string `json:"room_num"`
}
