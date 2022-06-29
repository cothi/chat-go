package ui

import (
	"time"

	"github.com/cothi/chat-go/client"
	"github.com/cothi/chat-go/utils"
	"github.com/marcusolsson/tui-go"
	ui "github.com/marcusolsson/tui-go"
)

type Cui struct {
	Ui      ui.UI
	Sidebar ui.Box
	History ui.Box
}

var cui Cui

// Append check box
func CheckboxRead(his *ui.Box, client *client.Client) {
	for {
		time.Sleep(1 * time.Second)
		m := <-client.Inbox
		cui.Ui.Update(func() {
			his.Append(ui.NewHBox(
				ui.NewLabel(time.Now().Format("15:14")),
				ui.NewLabel(m),
				ui.NewSpacer(),
			))
		})
	}
}

// ui setup
func UiSetup() {

	sidebar := ui.NewVBox(
		tui.NewLabel("CHANNELS"),
		tui.NewLabel("general"),
		tui.NewSpacer(),
	)

	sidebar.Insert(1, tui.NewLabel("check"))
	sidebar.SetBorder(true)
	sidebar.SetFocused(true)

	history := tui.NewVBox()
	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(ui.Expanding, ui.Expanding)

	input.OnSubmit(func(e *ui.Entry) {

		client.Outbox <- []byte(e.Text())
		input.SetText("")
	})

	root := ui.NewHBox(sidebar, chat)
	ui, err := ui.New(root)
	utils.Error_check(err)
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	cui.Ui = ui

	go CheckboxRead(history, client)
  cui.Ui.Run()

}

/* func StartClient(serverPort string) {
	client := client.ClientInit(serverPort)
	UiSetup(client)
} */
