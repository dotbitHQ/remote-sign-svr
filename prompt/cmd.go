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
		Label:    "请继续选择功能👇",
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
	t.cmdFunc["1.激活签名服务"] = func() error {
		return t.initRemoteSignSvr()
	}
	t.cmdFunc["2.创建钱包"] = func() error {
		return t.createWallet()
	}
	t.cmdFunc["3.导入钱包"] = func() error {
		return t.importWallet()
	}
	t.cmdFunc["4.获取钱包信息"] = func() error {
		return t.getWalletInfo()
	}
	t.cmdFunc["5.加密"] = func() error {
		return t.encData()
	}
	t.cmdFunc["6.解密"] = func() error {
		return t.decData()
	}
	t.cmdFunc["6.退出"] = func() error {
		os.Exit(1)
		return nil
	}
}
