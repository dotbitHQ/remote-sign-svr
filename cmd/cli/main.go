package main

import (
	"fmt"
	"remote-sign-svr/prompt"
)

func main() {
	var pt prompt.ToolPrompt

	if err := pt.InitRemoteSignSvr(); err != nil {
		fmt.Sprintln("Failed to int remote sign svr, err:", err.Error())
		return
	}

	if err := pt.Menu(); err != nil {
		fmt.Sprintln("Failed to show menu, err:", err.Error())
		return
	}
}
