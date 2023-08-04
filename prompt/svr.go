package prompt

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/manifoldco/promptui"
	"remote-sign-svr/http_svr/handle"
	"strings"
)

func (t *ToolPrompt) initKey() error {
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return fmt.Errorf("❌ The key cannot be empty")
		}
		if len(input) < KeyLen {
			return fmt.Errorf("❌ Key length is %d characters", KeyLen)
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Please enter the encryption key (%d characters)", KeyLen),
		Validate: validate,
		Mask:     '*',
	}
	key, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}

	validate2 := func(input string) error {
		if key != input {
			return fmt.Errorf("❌ The two keys are not the same")
		}
		return nil
	}
	prompt2 := promptui.Prompt{
		Label:    "Please enter the encrypted key again",
		Validate: validate2,
		Mask:     '*',
	}
	_, err = prompt2.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}

	t.key = key
	return nil
}

func (t *ToolPrompt) activateRemoteSignSvr() error {
	if t.key == "" {
		if err := t.initKey(); err != nil {
			return fmt.Errorf("initKey err: %s", err.Error())
		}
	}

	prompt := promptui.Prompt{
		Label: "Whether to activate the remote sign svr(y/n)",
	}
	key, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	if key != "y" && key != "Y" {
		return nil
	}

	url := fmt.Sprintf("http://%s/v1/init/svr", t.remoteSignSvr)
	req := handle.ReqInitSvr{Key: t.key}
	resp := handle.RespInitSvr{}
	if err := http_api.SendReq(url, req, &resp); err != nil {
		doErr(err, "❌ Failed to activate remote sign svr")
		return nil
	}
	return nil
}

func (t *ToolPrompt) getWalletInfo() error {
	prompt := promptui.Prompt{
		Label: "Please enter your wallet address",
	}
	addr, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	url := fmt.Sprintf("http://%s/v1/address/info", t.remoteSignSvr)
	req := handle.ReqAddressInfo{Address: addr}
	resp := handle.RespAddressInfo{}
	if err := http_api.SendReq(url, req, &resp); err != nil {
		doErr(err, "❌ Failed to get address info")
		return nil
	}
	msg := fmt.Sprintf(`\nAddress: %s
AddrChain: %s
Private: %s
CompressType: %t
AddrStatus: %d
Remark: %s
`, resp.Address, resp.AddrChain, resp.Private, resp.CompressType.Bool(), resp.AddrStatus, resp.Remark)
	fmt.Println(msg)
	return nil
}

func doErr(err error, msg string) {
	if strings.Contains(err.Error(), "connection refused") {
		fmt.Println("❌ Remote Sign Svr not started")
	} else {
		fmt.Println(msg, ", err: ", err.Error())
	}
}
