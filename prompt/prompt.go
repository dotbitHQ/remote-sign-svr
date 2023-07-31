package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strings"
)

type ToolPrompt struct {
	key        string
	cmdFunc    map[string]func() error
	walletFunc map[string]func() error
}

const (
	KeyLen = 16
)

func (t *ToolPrompt) InitKey() error {
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return fmt.Errorf("密钥不能为空！")
		}
		if len(input) < KeyLen {
			return fmt.Errorf("密钥长度不小于 %d 个字符", KeyLen)
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("请输入加密的密钥(%d个字符)", KeyLen),
		Validate: validate,
		Mask:     '*',
	}
	result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	t.key = result

	return t.checkKey()
}

func (t *ToolPrompt) checkKey() error {
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return fmt.Errorf("密钥不能为空！")
		}
		if len(input) < KeyLen {
			return fmt.Errorf("密钥长度不小于 %d 个字符", KeyLen)
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "请再次输入加密的密钥",
		Validate: validate,
		Mask:     '*',
	}
	result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	if t.key != result {
		return fmt.Errorf("两次密钥不一致")
	}
	return nil
}
