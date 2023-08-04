package prompt

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/manifoldco/promptui"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"remote-sign-svr/encrypt"
	"remote-sign-svr/http_svr/handle"
	"remote-sign-svr/tables"
	"remote-sign-svr/wallet"
	"sort"
)

func (t *ToolPrompt) createWallet() error {
	t.initWalletFunc()
	prompts := promptui.Select{
		HideHelp: true,
		Size:     len(t.walletFunc),
		Label:    "Please select the type of wallet you want to create 👇",
		Items:    t.getWalletNames(),
	}
	_, result, err := prompts.Run()
	if err != nil {
		return fmt.Errorf("prompts.Run() err: %s", err.Error())
	}
	if err := t.walletFunc[result](); err != nil {
		e := fmt.Errorf("❌ Failed to execute t.walletFunc[%s]() err:\n%s", result, err.Error())
		fmt.Println(e.Error())
	}
	return t.createWallet()
}

func (t *ToolPrompt) importNow(wa *wallet.AddressInfo) error {
	fmt.Println("Address:", wa.Address)
	fmt.Println("Private:", wa.Private)
	switch wa.AddrChain {
	case tables.AddrChainDOGE:
		fmt.Println("Payload:", wa.Payload)
		fmt.Println("Wif:", wa.WifStr)
	case tables.AddrChainCKB:
		fmt.Println("LockArgs:", wa.LockArgs)
	}
	//
	prompt := promptui.Prompt{
		Label: "Whether to import the wallet address(y/n)",
	}
	key, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	if key != "y" {
		return nil
	}
	if t.key == "" {
		fmt.Println("❌ Encryption key is empty")
		return t.Menu()
	}
	//
	prompt = promptui.Prompt{
		Label: "Please enter a note",
	}
	remark, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	//
	private, err := encrypt.AesEncrypt(wa.Private, t.key)
	if err != nil {
		return fmt.Errorf("AesEncrypt err: %s", err.Error())
	}

	url := fmt.Sprintf("http://%s/v1/import/address", t.remoteSignSvr)
	req := handle.ReqImportAddress{
		AddrChain: wa.AddrChain,
		Address:   wa.Address,
		Private:   private,
		Remark:    remark,
	}
	resp := handle.RespImportAddress{}
	if err := http_api.SendReq(url, req, &resp); err != nil {
		doErr(err, "❌ Failed to import address")
		return nil
	}
	return nil
}

func (t *ToolPrompt) initWalletFunc() {
	if t.walletFunc != nil {
		return
	}
	t.walletFunc = make(map[string]func() error)
	t.walletFunc["1.EVM"] = func() error {
		res, err := wallet.CreateWalletEVM()
		if err != nil {
			return fmt.Errorf("CreateWalletEVM err: %s", err.Error())
		}
		return t.importNow(res)
	}
	t.walletFunc["2.TRON"] = func() error {
		res, err := wallet.CreateWalletTRON()
		if err != nil {
			return fmt.Errorf("CreateWalletTRON err: %s", err.Error())
		}
		return t.importNow(res)
	}
	t.walletFunc["3.DOGE"] = func() error {
		res, err := wallet.CreateWalletDOGE(true)
		if err != nil {
			return fmt.Errorf("CreateWalletDOGE err: %s", err.Error())
		}
		return t.importNow(res)
	}
	t.walletFunc["4.CKB"] = func() error {
		prompt := promptui.Prompt{
			Label: "Whether to create a test network address(y/n)",
		}
		key, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("prompt.Run() err: %s", err.Error())
		}
		mode := address.Mainnet
		if key == "y" {
			mode = address.Testnet
		}
		res, err := wallet.CreateWalletCKB(mode)
		if err != nil {
			return fmt.Errorf("CreateWalletCKB err: %s", err.Error())
		}
		return t.importNow(res)
	}
	t.walletFunc["5.back"] = func() error {
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
	prompts := promptui.Select{
		HideHelp: true,
		Size:     len(t.walletFunc),
		Label:    "Please select the type of wallet to import 👇",
		Items:    []string{"1.EVM", "2.TRON", "3.DOGE", "4.CKB", "5.back"},
	}
	index, _, err := prompts.Run()
	if err != nil {
		return fmt.Errorf("prompts.Run() err: %s", err.Error())
	}
	var addrChain tables.AddrChain
	switch index {
	case 0:
		addrChain = tables.AddrChainEVM
	case 1:
		addrChain = tables.AddrChainTRON
	case 2:
		addrChain = tables.AddrChainDOGE
	case 3:
		addrChain = tables.AddrChainCKB
	default:
		return nil
	}

	//
	prompt := promptui.Prompt{
		Label: "Please enter your wallet address",
	}
	addr, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	//
	prompt = promptui.Prompt{
		Label: "Please enter the encryption private key",
	}
	private, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	//
	prompt = promptui.Prompt{
		Label: "Please enter a note",
	}
	remark, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt.Run() err: %s", err.Error())
	}
	//
	url := fmt.Sprintf("http://%s/v1/import/address", t.remoteSignSvr)
	req := handle.ReqImportAddress{
		AddrChain: addrChain,
		Address:   addr,
		Private:   private,
		Remark:    remark,
	}
	resp := handle.RespImportAddress{}
	if err := http_api.SendReq(url, req, &resp); err != nil {
		doErr(err, "❌ Failed to import address")
		return nil
	}
	return nil
}
