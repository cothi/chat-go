package main

import (
	"fmt"
	"strconv"
	"time"

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
	chat string
	time string
}

type Tui struct {
	ui tui.UI
}

var rootUi Tui

func (room *ChatRoom) Init() {
	post := &Post{
		chat: "t",
		time: time.Now().Local().Format("15:04"),
	}

	post2 := &Post{
		chat: "t2",
		time: time.Now().Local().Format("15:04"),
	}

	chatP := &ChatPost{
		ChatHistory: []*Post{post},
	}
	chatT := &ChatPost{
		ChatHistory: []*Post{post2},
	}

	fmt.Println(chatP, chatT)

	room.ChatRoom = make(map[string]*ChatPost)
	room.ChatRoom["0"] = chatP
	room.ChatRoom["1"] = chatT
}

func historyUpdate(ui *tui.Box, selected string) {

	// b := *ChatRooms.ChatRoom[selected].ChatHistory[0]

	// fmt.Println(b.chat)

	newHbox := tui.NewHBox()

	newHbox.SetBorder(true)
	for _, m := range ChatRooms.ChatRoom[selected].ChatHistory {
		// fmt.Println(*&m.chat)
		newHbox.Append(
			tui.NewHBox(
				tui.NewLabel(*&m.time),
				tui.NewLabel(*&m.chat),
			))
	}
	*ui = *newHbox

	// c := ChatRooms.ChatRoom[selected].ChatHistory
	// fmt.Println(c)

	// for _, m := range ChatRooms.ChatRoom[selected].ChatHistory {
	// }
}

// ui setup
func main() {
	ChatRooms.Init()

	barList := tui.NewList()
	barList.AddItems(
		"Test",
		"Test2",
	)

	sidebar := tui.NewHBox(
		barList,
	)
	sidebar.SetBorder(true)

	chatHistory := tui.NewVBox()
	chatHistory.SetBorder(true)

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
	rootUi.ui = root

	utils.Error_check(err)
	barList.OnSelectionChanged(
		func(ui *tui.List) {
			// fmt.Println(tui.Selected())
			historyUpdate(chatHistory, strconv.Itoa(ui.Selected()))
			/*
				*chatHistory = *tui.NewVBox(
					tui.NewLabel(ui.SelectedItem())b,
				) */
			// root.SetWidget(entirePanel)
		},
	)

	root.SetKeybinding("Esc", func() { root.Quit() })

	root.Run()
}
