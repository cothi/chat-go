package ui

import (
	"encoding/json"
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
	ChatHistory []*client.Post
}

type Tui struct {
	Ui         tui.UI
	HistoryBox *tui.Box
	NowRoom    string
	ChanList   *tui.List
}

var rootUi Tui
var cli client.Client

// Init method
func (room *ChatRoom) Init() {

	cli.Name = "anonymous"
	rootUi.NowRoom = "Lobby"

	post := &client.Post{
		Chat:     "t",
		Time:     time.Now().Local().Format("15:04"),
		Roomname: "Lobby",
		Kind:     client.MessageCreateMsg,
	}

	post2 := &client.Post{
		Chat:     "t2",
		Time:     time.Now().Local().Format("15:04"),
		Roomname: "tmp",
		Kind:     client.MessageCreateMsg,
	}

	chatP := &ChatPost{
		ChatHistory: []*client.Post{post},
	}
	chatT := &ChatPost{
		ChatHistory: []*client.Post{post2},
	}

	room.ChatRoom = make(map[string]*ChatPost)
	room.ChatRoom["Lobby"] = chatP
	room.ChatRoom["tmp"] = chatT
	for {
		var postMsg client.Post
		msg := <-cli.Inbox
		json.Unmarshal(msg, &postMsg)
		room.HandleRecievedMessage(&postMsg)

		/* rootUi.Ui.Update(func() {
			ui.Append(
				tui.NewHBox(
					tui.NewLabel(postMsg.Chat+" "),
					tui.NewLabel(postMsg.Time+" : "),
					tui.NewLabel(postMsg.Chat),
					tui.NewSpacer(),
				),
			)
		}) */
	}
}

func (room *ChatRoom) HandleRecievedMessage(msg *client.Post) {
	switch msg.Kind {
	case client.MessageCreateMsg:
		room.AddMessage(msg)
	case client.MessageCreateRoom:
		room.CreateRoom(msg)
	case client.MessageJoinRoom:
		room.AddMessage(msg)
	}
}

// AddMessage function
func (room *ChatRoom) AddMessage(p *client.Post) {
	room.ChatRoom[p.Roomname].ChatHistory = append(room.ChatRoom[p.Roomname].ChatHistory, p)
	// fmt.Println(room.ChatRoom[p.Roomname].ChatHistory)
	room.HistoryUpdate(p.Roomname)
}

func (room *ChatRoom) CreateRoom(p *client.Post) {
	if _, ok := room.ChatRoom[p.Chat]; ok {
		return
	}
	msg := &*p
	tmp := &ChatPost{
		ChatHistory: []*client.Post{msg},
	}
	room.ChatRoom[p.Chat] = tmp

	rootUi.Ui.Update(func() {
		rootUi.ChanList.AddItems(p.Chat)
	})
}

// HistoryUpdate function
func (room *ChatRoom) HistoryUpdate(selected string) {
	newVbox := tui.NewVBox()
	newVbox.Append(tui.NewSpacer())

	for _, m := range room.ChatRoom[selected].ChatHistory {
		msg := tui.NewHBox(
			tui.NewLabel(m.Nickname+" "),
			tui.NewLabel(m.Time+" : "),
			tui.NewLabel(m.Chat),
			tui.NewSpacer(),
		)
		// fmt.Println(m.Chat)
		newVbox.Append(msg)
	}

	go rootUi.Ui.Update(func() {
		*rootUi.HistoryBox = *newVbox
	})
}

// Ui Setup
func UiSetup(serverPort string) {
	cli.ClientInit(serverPort)

	barList := tui.NewList()
	barList.AddItems(
		"Lobby",
		"tmp",
	)

	sidebar := tui.NewVBox(
		tui.NewLabel("CHANNEL"),
		barList,
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	chatHistory := tui.NewVBox()

	chatScroll := tui.NewScrollArea(chatHistory)
	chatScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(chatScroll)
	historyBox.SetBorder(true)

	chatEntry := tui.NewEntry()
	chatEntry.SetFocused(true)

	chatInputBox := tui.NewHBox(chatEntry)
	chatInputBox.SetBorder(true)
	chatInputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chatPanel := tui.NewVBox(historyBox, chatInputBox)
	chatPanel.SetSizePolicy(tui.Expanding, tui.Maximum)

	entirePanel := tui.NewHBox(sidebar, chatPanel)
	entirePanel.SetSizePolicy(tui.Expanding, tui.Maximum)

	tui.DefaultFocusChain.Set(barList, chatEntry)

	root, err := tui.New(entirePanel)
	rootUi.Ui = root

	utils.Error_check(err)
	barList.OnSelectionChanged(
		func(ui *tui.List) {
			rootUi.NowRoom = ui.SelectedItem()
			ChatRooms.HistoryUpdate(ui.SelectedItem())
		},
	)

	chatEntry.OnSubmit(func(e *tui.Entry) {
		p := client.Post{
			Chat:     e.Text(),
			Time:     time.Now().Local().Format("15:04"),
			Roomname: rootUi.NowRoom,
			Nickname: cli.Name,
			Kind:     client.MessageCreateMsg,
		}
		// pBytes, _ := json.Marshal(p)
		// cli.Outbox <- pBytes
		cli.SendHandleMessage(&p)
		e.SetText("")
	})

	rootUi.ChanList = barList
	rootUi.HistoryBox = chatHistory
	root.SetKeybinding("Esc", func() { root.Quit() })

	go ChatRooms.Init()
	root.Run()
}
