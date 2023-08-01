package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"sort"
)

func (t *ToolPrompt) Menu() error {
	t.initCmdFunc()
	prompts := promptui.Select{
		HideHelp: true,
		Size:     len(t.cmdFunc),
		Label:    "Please continue to select function 👇",
		Items:    t.getCmdNames(),
	}
	_, result, err := prompts.Run()
	if err != nil {
		return fmt.Errorf("prompts.Run() err: %s", err.Error())
	}
	if err := t.cmdFunc[result](); err != nil {
		return fmt.Errorf("p.cmdFunc[%s]() err: %s", result, err.Error())
	}
	return t.Menu()
}

func (t *ToolPrompt) getCmdNames() []string {
	var names []string
	for name, _ := range t.cmdFunc {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (t *ToolPrompt) initCmdFunc() {
	if t.cmdFunc != nil {
		return
	}
	t.cmdFunc = make(map[string]func() error)
	t.cmdFunc["1.Activate Remote Sign Svr"] = func() error {
		return t.initRemoteSignSvr()
	}
	t.cmdFunc["2.Create Wallet"] = func() error {
		return t.createWallet()
	}
	t.cmdFunc["3.Import Wallet"] = func() error {
		return t.importWallet()
	}
	t.cmdFunc["4.Search Wallet Info"] = func() error {
		return t.getWalletInfo()
	}
	t.cmdFunc["5.Encrypted Data"] = func() error {
		return t.encData()
	}
	t.cmdFunc["6.Decrypted Data"] = func() error {
		return t.decData()
	}
	t.cmdFunc["6.Exit"] = func() error {
		os.Exit(1)
		return nil
	}
}
