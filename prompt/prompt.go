package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

type ToolPrompt struct {
	remoteSignSvr string
	key           string // 1234567890098765
	cmdFunc       map[string]func() error
	walletFunc    map[string]func() error
}

const (
	KeyLen = 16
)

func (t *ToolPrompt) InitRemoteSignSvr() error {
	prompt := promptui.Prompt{
		Label: "Please enter the IP and port of remote sign svr(127.0.0.1:9093)",
	}
	svr, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err:\n%s", err.Error())
	}
	t.remoteSignSvr = svr
	if t.remoteSignSvr == "" {
		t.remoteSignSvr = "127.0.0.1:9093"
	}
	return nil
}
