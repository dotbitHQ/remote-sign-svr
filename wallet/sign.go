package wallet

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/sign"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/secp256k1"
	"math/big"
	"remote-sign-svr/tables"
)

var (
	ErrUnsupportedAddrChain = errors.New("unsupported address chain")
	ErrUnsupportedSignType  = errors.New("unsupported sign type")
)

type SignType int

const (
	SignTypeTx  SignType = 0
	SignTypeMsg SignType = 1
)

func Sign(signType SignType, addrInfo tables.TableAddressInfo, data string, chainId *big.Int) (string, error) {
	switch signType {
	case SignTypeTx:
		return signTx(addrInfo, data, chainId)
	case SignTypeMsg:
		return signMsg(addrInfo, data)
	default:
		return "", ErrUnsupportedSignType
	}
}

func signTx(addrInfo tables.TableAddressInfo, data string, chainId *big.Int) (string, error) {
	private := addrInfo.Private
	switch addrInfo.AddrChain {
	case tables.AddrChainEVM: // tx hex
		bys, err := hex.DecodeString(data)
		if err != nil {
			return "", fmt.Errorf("hex.DecodeString err: %s", err.Error())
		}
		tx := &types.Transaction{}
		if err = rlp.DecodeBytes(bys, tx); err != nil {
			return "", fmt.Errorf("rlp.DecodeBytes err: %s", err.Error())
		}
		privateKey, err := crypto.HexToECDSA(private)
		if err != nil {
			return "", fmt.Errorf("crypto.HexToECDSA err: %s", err.Error())
		}
		sigTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
		if err != nil {
			return "", fmt.Errorf("SignTx err: %s", err.Error())
		}
		sigData, err := rlp.EncodeToBytes(sigTx)
		if err != nil {
			return "", fmt.Errorf("rlp.EncodeToBytes err: %s", err.Error())
		}
		return hex.EncodeToString(sigData), nil
	case tables.AddrChainTRON: // sig hex
		bys, err := hex.DecodeString(data)
		if err != nil {
			return "", fmt.Errorf("hex.DecodeString err: %s", err.Error())
		}
		privateKey, err := crypto.HexToECDSA(private)
		if err != nil {
			return "", fmt.Errorf("crypto.HexToECDSA err: %s", err.Error())
		}
		sig, err := crypto.Sign(bys, privateKey)
		if err != nil {
			return "", fmt.Errorf("crypto.Sign err: %s", err.Error())
		}
		return hex.EncodeToString(sig), nil
	case tables.AddrChainDOGE: // tx hex
		bys, err := hex.DecodeString(data)
		if err != nil {
			return "", fmt.Errorf("hex.DecodeString err: %s", err.Error())
		}
		var tx wire.MsgTx
		if err = tx.DeserializeNoWitness(bytes.NewReader(bys)); err != nil {
			return "", fmt.Errorf("tx.DeserializeNoWitness err: %s", err.Error())
		}
		for i := 0; i < len(tx.TxIn); i++ {
			pkScript, privateKey, compress, err := bitcoin.HexPrivateKeyToScript(addrInfo.Address, bitcoin.GetDogeMainNetParams(), private)
			if err != nil {
				return "", fmt.Errorf("HexPrivateKeyToScript err: %s", err.Error())
			}
			sig, err := txscript.SignatureScript(&tx, i, pkScript, txscript.SigHashAll, privateKey, compress)
			if err != nil {
				return "", fmt.Errorf("SignatureScript err: %s", err.Error())
			}
			tx.TxIn[i].SignatureScript = sig
		}
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSizeStripped()))
		_ = tx.SerializeNoWitness(buf)
		return hex.EncodeToString(buf.Bytes()), nil
	case tables.AddrChainCKB: // sig hex
		bys, err := hex.DecodeString(data)
		if err != nil {
			return "", fmt.Errorf("hex.DecodeString err: %s", err.Error())
		}
		key, err := secp256k1.HexToKey(private)
		if err != nil {
			return "", fmt.Errorf("secp256k1.HexToKey err: %s", err.Error())
		}
		sig, err := key.Sign(bys)
		if err != nil {
			return "", fmt.Errorf("key.Sign err: %s", err.Error())
		}
		return hex.EncodeToString(sig), nil
	default:
		return "", ErrUnsupportedAddrChain
	}
}

func signMsg(addrInfo tables.TableAddressInfo, data string) (string, error) {
	private := addrInfo.Private
	switch addrInfo.AddrChain {
	case tables.AddrChainEVM:
		sig, err := sign.PersonalSignature([]byte(data), private)
		if err != nil {
			return "", fmt.Errorf("sign.PersonalSignature err: %s", err.Error())
		}
		return hex.EncodeToString(sig), nil
	case tables.AddrChainTRON:
		sig, err := sign.TronSignature(true, []byte(data), private)
		if err != nil {
			return "", fmt.Errorf("sign.TronSignature err: %s", err.Error())
		}
		return hex.EncodeToString(sig), nil
	case tables.AddrChainDOGE:
		sig, err := sign.DogeSignature([]byte(data), private, addrInfo.CompressType.Bool())
		if err != nil {
			return "", fmt.Errorf("sign.DogeSignature err: %s", err.Error())
		}
		return hex.EncodeToString(sig), nil
	default:
		return "", ErrUnsupportedAddrChain
	}
}
