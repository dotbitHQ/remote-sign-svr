package wallet

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"strings"
)

func CreateWalletEVM() error {
	key, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("GenerateKey err: %s", err.Error())
	}
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(crypto.FromECDSA(key))

	fmt.Println("钱包地址:", address)
	fmt.Println("私钥:", privateKey)
	return nil
}

func CreateWalletTRON() error {
	key, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("GenerateKey err: %s", err.Error())
	}
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(crypto.FromECDSA(key))

	address = "41" + address[2:]
	address, err = common.TronHexToBase58(address)
	if err != nil {
		return fmt.Errorf("TronHexToBase58 err: %s", err.Error())
	}

	fmt.Println("钱包地址:", address)
	fmt.Println("私钥:", privateKey)
	return nil
}

func CreateWalletDOGE(compress bool) error {
	mainNetParams := bitcoin.GetDogeMainNetParams()
	key, err := btcec.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("NewPrivateKey err: %s", err.Error())
	}
	wif, err := btcutil.NewWIF(key, &mainNetParams, compress)
	if err != nil {
		return fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
	}
	addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &mainNetParams)
	if err != nil {
		return fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
	}

	payload := hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:])
	wifStr := wif.String()
	address := addressPubKey.EncodeAddress()
	privateKey := hex.EncodeToString(key.Serialize())

	fmt.Println("Payload:", payload)
	fmt.Println("WIF:", wifStr)
	fmt.Println("钱包地址:", address)
	fmt.Println("私钥:", privateKey)
	return nil
}

func CreateWalletCKB(mode address.Mode) error {
	res, err := address.GenerateShortAddress(mode)
	if err != nil {
		return fmt.Errorf("GenerateAddress err: %s", err.Error())
	}
	fmt.Println("args:", res.LockArgs)
	fmt.Println("钱包地址:", res.Address)
	fmt.Println("私钥:", strings.TrimPrefix(res.PrivateKey, "0x"))
	return nil
}
