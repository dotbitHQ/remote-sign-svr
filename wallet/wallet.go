package wallet

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"remote-sign-svr/tables"
	"strings"
)

type AddressInfo struct {
	AddrChain tables.AddrChain `json:"addr_chain"`
	Address   string           `json:"address"`
	Private   string           `json:"private"`
	Payload   string           `json:"payload"`
	WifStr    string           `json:"wif_str"`
	LockArgs  string           `json:"lock_args"`
}

func CreateWalletEVM() (*AddressInfo, error) {

	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("GenerateKey err: %s", err.Error())
	}
	addr := crypto.PubkeyToAddress(key.PublicKey).Hex()
	private := hex.EncodeToString(crypto.FromECDSA(key))
	return &AddressInfo{AddrChain: tables.AddrChainEVM, Address: addr, Private: private}, nil
}

func CreateWalletTRON() (*AddressInfo, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("GenerateKey err: %s", err.Error())
	}
	addr := crypto.PubkeyToAddress(key.PublicKey).Hex()
	private := hex.EncodeToString(crypto.FromECDSA(key))

	addr = "41" + addr[2:]
	addr, err = common.TronHexToBase58(addr)
	if err != nil {
		return nil, fmt.Errorf("TronHexToBase58 err: %s", err.Error())
	}
	return &AddressInfo{AddrChain: tables.AddrChainTRON, Address: addr, Private: private}, nil
}

func CreateWalletDOGE(compress bool) (*AddressInfo, error) {
	mainNetParams := bitcoin.GetDogeMainNetParams()
	key, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("NewPrivateKey err: %s", err.Error())
	}
	wif, err := btcutil.NewWIF(key, &mainNetParams, compress)
	if err != nil {
		return nil, fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
	}
	addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &mainNetParams)
	if err != nil {
		return nil, fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
	}

	payload := hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:])
	wifStr := wif.String()
	addr := addressPubKey.EncodeAddress()
	private := hex.EncodeToString(key.Serialize())

	return &AddressInfo{AddrChain: tables.AddrChainDOGE, Address: addr, Private: private, Payload: payload, WifStr: wifStr}, nil
}

func CreateWalletCKB(mode address.Mode) (*AddressInfo, error) {
	res, err := address.GenerateShortAddress(mode)
	if err != nil {
		return nil, fmt.Errorf("GenerateAddress err: %s", err.Error())
	}
	return &AddressInfo{AddrChain: tables.AddrChainCKB, Address: res.Address, Private: strings.TrimPrefix(res.PrivateKey, "0x"), LockArgs: res.LockArgs}, nil
}

func CreateWalletBTC(netParams chaincfg.Params) (*AddressInfo, error) {
	key, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("NewPrivateKey err: %s", err.Error())
	}
	wif, err := btcutil.NewWIF(key, &netParams, true)
	if err != nil {
		return nil, fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
	}
	addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &netParams)
	if err != nil {
		return nil, fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
	}
	pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]

	addressWPH, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, &netParams)
	if err != nil {
		return nil, fmt.Errorf("NewAddressWitnessPubKeyHash err: %s", err.Error())
	}

	payload := hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:])
	wifStr := wif.String()
	addr := addressWPH.EncodeAddress()
	private := hex.EncodeToString(key.Serialize())

	return &AddressInfo{AddrChain: tables.AddrChainDOGE, Address: addr, Private: private, Payload: payload, WifStr: wifStr}, nil
}
