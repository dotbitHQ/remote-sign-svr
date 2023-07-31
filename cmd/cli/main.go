package main

import (
	"github.com/scorpiotzh/mylog"
	"remote-sign-svr/prompt"
)

var log = mylog.NewLogger("main", mylog.LevelDebug)

func main() {
	var pt prompt.ToolPrompt
	//if err := pt.InitKey(); err != nil {
	//	log.Errorf("pt.InitKey() err: %s", err.Error())
	//}
	if err := pt.Menu(); err != nil {
		log.Error(err.Error())
	}
}
