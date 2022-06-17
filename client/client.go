package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/cothi/tcp-chat-remodel/utils"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Client struct {
	conn net.Conn
	m    sync.Mutex
}

var (
	wg     sync.WaitGroup
	client Client
)

func (c *Client) Read() {

	recv := make([]byte, 4096)

	for {
		n, e := c.conn.Read(recv)
		utils.Error_check(e)

		str := strings.TrimSpace(string(recv[:n]))
		fmt.Printf("recv: %s\n\n", str)
	}
}

func (c *Client) Write() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		c.conn.Write([]byte(text))
	}

}

func client_ui() {

	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Text = "Hello world"
	p.SetRect(0, 0, 100, 100)
	ui.Render()

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}

// Reads from Stdin, and outputs to the socket.
func main() {
  client_ui()
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	client.conn = conn

	go client.Write()
	go client.Read()

	utils.Error_check(err)

	wg.Add(1)
	wg.Wait()

}
