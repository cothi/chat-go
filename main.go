package main

import "github.com/cothi/tcp-chat-remodel/cmd"

func main() {
	/* var wg sync.WaitGroup

	wg.Add(1)

	// var cui ui.Cui
	client := client.ClientInit()
	// ui.UiSetup(&cui, client)

	// go cui.CheckboxRead(client)
	// go cui.Ui.Run()

	go ui.UiSetup(client)
	wg.Wait() */

  cmd.Execute()
}

