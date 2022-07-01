package server

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/cothi/chat-go/utils"
)

type MessageKind int

const (
	MessageCreateMsg MessageKind = iota
	MessageCreateRoom
	MessageJoinRoom
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

type Post struct {
	Chat     string `json:"chat"`
	Time     string `json:"time"`
	Roomname string `json:"roomname"`
	NickName string `json:"nickname"`
	Kind     int    `json:"kind"`
}

func (c *Chan) Init() {
	c.room = make(map[string]*room)
	room1 := room{
		room_name: "Lobby",
	}
	c.room["Lobby"] = &room1
}

func (c *Chan) Create(p *Post, client *Client) {
	roomName := p.Chat
	c.room[p.Chat] = &room{
		room_name: roomName,
	}
	c.room[roomName].clients = append(c.room[roomName].clients, client)
	c.room[roomName].message = []string{"create room"}
	fmt.Println("create room", roomName)

	b, e := json.Marshal(p)
	utils.Error_check(e)

	for _, r := range c.room {
		r.broadcast(b)
	}
}

func (c *Chan) Join(roomName string, client *Client) {
	c.room[roomName].clients = append(c.room[roomName].clients, client)
	c.room[roomName].message = append(c.room[roomName].message, "join")
	fmt.Println("join room", client.name)
}

type room struct {
	room_name string
	clients   []*Client
	message   []string
}

func (r *room) broadcast(msg []byte) {
	for _, c := range r.clients {
		fmt.Println(msg)
		c.box <- msg
	}
}

type Client struct {
	name string
	conn net.Conn
	box  chan []byte
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
		fmt.Println(msg)
		HandleMessage(&msg, c)
	}
}

func InitClient(conn net.Conn) {
	client := Client{
		name: "anonymous",
		conn: conn,
		box:  make(chan []byte),
	}
	Lobby.room["Lobby"].clients = append(Lobby.room["Lobby"].clients, &client)
	go client.Write()
	go client.Read()
}

func HandleMessage(p *Post, c *Client) {
	switch p.Kind {
	case int(MessageCreateMsg):
		b, e := json.Marshal(p)
		utils.Error_check(e)
		Lobby.room[p.Roomname].broadcast(b)
	case int(MessageCreateRoom):
		Lobby.Create(p, c)
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
