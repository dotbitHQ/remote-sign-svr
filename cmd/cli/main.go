package main

import (
	"fmt"
	"remote-sign-svr/prompt"
	"time"
)

func main() {
	var pt prompt.ToolPrompt

	if err := pt.InitRemoteSignSvr(); err != nil {
		fmt.Println("❌ Failed to int remote sign svr, err:\n", err.Error())
		return
	}

	if err := pt.Menu(); err != nil {
		fmt.Println("❌ Failed to show menu, err:\n", err.Error())
		return
	}
	time.Sleep(time.Second * 3)
}
