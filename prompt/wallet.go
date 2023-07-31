package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"remote-sign-svr/wallet"
	"sort"
)

func (t *ToolPrompt) createWallet() error {
	t.initWalletFunc()
	prompts := promptui.Select{
		HideHelp: true,
		Size:     len(t.walletFunc),
		Label:    "请选择要创建的钱包类型👇",
		Items:    t.getWalletNames(),
	}
	_, result, err := prompts.Run()
	if err != nil {
		return fmt.Errorf("prompts.Run() err: %s", err.Error())
	}
	if err := t.walletFunc[result](); err != nil {
		return fmt.Errorf("p.cmdFunc[%s]() err: %s", result, err.Error())
	}
	return t.createWallet()
}

func (t *ToolPrompt) initWalletFunc() {
	if t.walletFunc != nil {
		return
	}
	t.walletFunc = make(map[string]func() error)
	t.walletFunc["1.EVM"] = func() error {
		return wallet.CreateWalletEVM()
	}
	t.walletFunc["2.TRON"] = func() error {
		return wallet.CreateWalletTRON()
	}
	t.walletFunc["3.DOGE"] = func() error {
		return wallet.CreateWalletDOGE(true)
	}
	t.walletFunc["4.CKB"] = func() error {
		return wallet.CreateWalletCKB(address.Mainnet)
	}
	t.walletFunc["5.返回"] = func() error {
		return t.Menu()
	}
}

func (t *ToolPrompt) getWalletNames() []string {
	var names []string
	for name, _ := range t.walletFunc {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (t *ToolPrompt) importWallet() error {
	return nil
}
