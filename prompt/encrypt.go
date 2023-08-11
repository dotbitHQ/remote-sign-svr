package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"remote-sign-svr/encrypt"
	"strings"
)

func (t *ToolPrompt) encData() error {
	if t.key == "" {
		fmt.Println("❌ Encryption key is empty")
		return t.Menu()
	}
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return fmt.Errorf("❌ Encrypted data cannot be empty")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Please enter the data to be encrypted",
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
	fmt.Printf("✅ Encryption Result: \n%s\n", target)
	return nil
}

func (t *ToolPrompt) decData() error {
	if t.key == "" {
		fmt.Println("❌ Encryption key is empty")
		return t.Menu()
	}
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return fmt.Errorf("❌ Decrypted data cannot be empty")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Please enter the data to be decrypted",
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
	fmt.Printf("✅ Decryption Result: \n%s\n", target)
	return nil
}
