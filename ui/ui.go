package ui

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/cothi/chat-go/client"
	"github.com/cothi/chat-go/utils"
	"github.com/marcusolsson/tui-go"
)

type ChatRoom struct {
	ChatRoom map[string]*ChatPost
}

var ChatRooms ChatRoom

type ChatPost struct {
	ChatHistory []*Post
}

type Post struct {
	Chat    string `json:"chat"`
	Time    string `json:"time"`
	RoomNum string `json:"room_num"`
}

type Tui struct {
	Ui         tui.UI
	HistoryBox *tui.Box
	NowRoom    string
}

var rootUi Tui
var cli client.Client

// Init method
func (room *ChatRoom) Init() {

	rootUi.NowRoom = "0"
	cli.Name = "anonymous"

	post := &Post{
		Chat:    "t",
		Time:    time.Now().Local().Format("15:04"),
		RoomNum: "1",
	}

	post2 := &Post{
		Chat:    "t2",
		Time:    time.Now().Local().Format("15:04"),
		RoomNum: "2",
	}

	chatP := &ChatPost{
		ChatHistory: []*Post{post},
	}
	chatT := &ChatPost{
		ChatHistory: []*Post{post2},
	}

	room.ChatRoom = make(map[string]*ChatPost)
	room.ChatRoom["0"] = chatP
	room.ChatRoom["1"] = chatT
	var postMsg Post
	for {
		msg := <-cli.Inbox
		json.Unmarshal(msg, &postMsg)
		AddMessage(&postMsg)
	}
}

// AddMessage function
func AddMessage(p *Post) {

	*&ChatRooms.ChatRoom[p.RoomNum].ChatHistory = append(*&ChatRooms.ChatRoom[p.RoomNum].ChatHistory, p)
	HistoryUpdate(rootUi.HistoryBox, p.RoomNum)
}

// HistoryUpdate function
func HistoryUpdate(ui *tui.Box, selected string) {
	newVbox := tui.NewVBox()
	newVbox.SetBorder(true)
	newVbox.Append(tui.NewSpacer())

	rootUi.Ui.Update(func() {
		for _, m := range ChatRooms.ChatRoom[selected].ChatHistory {
			newVbox.Append(
				tui.NewHBox(
					tui.NewLabel(cli.Name+" "),
					tui.NewLabel(*&m.Time+" : "),
					tui.NewLabel(*&m.Chat),
					tui.NewSpacer(),
				))
		}
	})
	*ui = *newVbox
}

// Ui Setup
func UiSetup(serverPort string) {
	go ChatRooms.Init()
	cli.ClientInit(serverPort)

	barList := tui.NewList()
	barList.AddItems(
		"Test",
		"Test2",
	)

	sidebar := tui.NewVBox(
		tui.NewLabel("CHANNEL"),
		barList,
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	chatHistory := tui.NewVBox()
	chatHistory.SetBorder(true)

	chatScroll := tui.NewScrollArea(chatHistory)
	chatScroll.SetAutoscrollToBottom(true)

	rootUi.HistoryBox = chatHistory

	chatEntry := tui.NewEntry()
	chatInputBox := tui.NewHBox(chatEntry)
	chatInputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	chatInputBox.SetBorder(true)

	chatPanel := tui.NewVBox(chatHistory, chatInputBox)
	chatPanel.SetSizePolicy(tui.Expanding, tui.Maximum)

	entirePanel := tui.NewHBox(sidebar, chatPanel)
	entirePanel.SetSizePolicy(tui.Expanding, tui.Maximum)

	tui.DefaultFocusChain.Set(barList, chatEntry)

	root, err := tui.New(entirePanel)
	rootUi.Ui = root

	utils.Error_check(err)
	barList.OnSelectionChanged(
		func(ui *tui.List) {
			rootUi.NowRoom = strconv.Itoa(ui.Selected())
			HistoryUpdate(chatHistory, strconv.Itoa(ui.Selected()))
		},
	)

	chatEntry.OnSubmit(func(e *tui.Entry) {
		p := Post{
			Chat:    e.Text(),
			Time:    time.Now().Local().Format("15:04"),
			RoomNum: rootUi.NowRoom,
		}
		pBytes, _ := json.Marshal(p)
		cli.Outbox <- pBytes
		e.SetText("")
	})

	root.SetKeybinding("Esc", func() { root.Quit() })
	root.Run()
}
