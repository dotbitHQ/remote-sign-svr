package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"remote-sign-svr/encrypt"
	"strings"
)

func (t *ToolPrompt) encData() error {
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return fmt.Errorf("加密数据不能为空")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "请输入要加密的数据",
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	target, err := encrypt.AesEncrypt(result, t.key)
	if err != nil {
		return fmt.Errorf("AesEncrypt err: %s", err.Error())
	}
	fmt.Printf("加密结果: \n%s\n", string(target))
	return nil
}

func (t *ToolPrompt) decData() error {
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return fmt.Errorf("解密数据不能为空")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "请输入要解密的数据",
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	target, err := encrypt.AesDecrypt(result, t.key)
	if err != nil {
		return fmt.Errorf("AesDecrypt err: %s", err.Error())
	}
	fmt.Printf("解密结果: \n%s\n", string(target))
	return nil
}
