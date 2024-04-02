package wallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/sign"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/secp256k1"
	"math/big"
	"remote-sign-svr/tables"
	"strings"
)

var (
	ErrUnsupportedAddrChain = errors.New("unsupported address chain")
	ErrUnsupportedSignType  = errors.New("unsupported sign type")
)

type SignType int

const (
	SignTypeTx     SignType = 0
	SignTypeMsg    SignType = 1
	SignTypeETH712 SignType = 2
)

func Sign(signType SignType, addrInfo tables.TableAddressInfo, data string, chainId *big.Int, mmJsonObj *common.MMJsonObj) (string, error) {
	switch signType {
	case SignTypeTx:
		return signTx(addrInfo, data, chainId)
	case SignTypeMsg:
		return signMsg(addrInfo, data)
	case SignTypeETH712:
		return sign712(addrInfo, chainId.Int64(), data, mmJsonObj)
	default:
		return "", ErrUnsupportedSignType
	}
}

func signTx(addrInfo tables.TableAddressInfo, data string, chainId *big.Int) (string, error) {
	private := addrInfo.Private
	bys, err := hex.DecodeString(strings.TrimPrefix(data, "0x"))
	if err != nil {
		return "", fmt.Errorf("hex.DecodeString err: %s", err.Error())
	}

	switch addrInfo.AddrChain {
	case tables.AddrChainEVM: // tx hex
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
	case tables.AddrChainBTC:
		var tx wire.MsgTx
		if err = tx.DeserializeNoWitness(bytes.NewReader(bys)); err != nil {
			return "", fmt.Errorf("tx.DeserializeNoWitness err: %s", err.Error())
		}
		for i := 0; i < len(tx.TxIn); i++ {
			pkScript, privateKey, compress, err := bitcoin.HexPrivateKeyToScript(addrInfo.Address, bitcoin.GetBTCMainNetParams(), private)
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
	case tables.AddrChainBTC:
		segwitType := sign.P2PKH
		addrType, _, _ := bitcoin.FormatBTCAddr(addrInfo.Address)
		switch addrType {
		case bitcoin.BtcAddressTypeP2WPKH:
			segwitType = sign.P2WPKH
		case bitcoin.BtcAddressTypeP2PKH:
			segwitType = sign.P2PKH
		}
		sig, err := sign.BitcoinSignature([]byte(data), private, addrInfo.CompressType.Bool(), segwitType)
		if err != nil {
			return "", fmt.Errorf("sign.DogeSignature err: %s", err.Error())
		}
		return hex.EncodeToString(sig), nil
	default:
		return "", ErrUnsupportedAddrChain
	}
}

func sign712(addrInfo tables.TableAddressInfo, chainId int64, signMsg string, mmJsonObj *common.MMJsonObj) (string, error) {
	log.Info("sign712:", chainId, signMsg, addrInfo.Address)
	var signData []byte
	private := addrInfo.Private

	var obj3 apitypes.TypedData
	mmJson := mmJsonObj.String()
	oldChainId := fmt.Sprintf("chainId\":%d", chainId)
	newChainId := fmt.Sprintf("chainId\":\"%d\"", chainId)
	mmJson = strings.ReplaceAll(mmJson, oldChainId, newChainId)
	oldDigest := "\"digest\":\"\""
	newDigest := fmt.Sprintf("\"digest\":\"%s\"", signMsg)
	mmJson = strings.ReplaceAll(mmJson, oldDigest, newDigest)

	_ = json.Unmarshal([]byte(mmJson), &obj3)
	var mmHash, signature []byte
	mmHash, signature, err := sign.EIP712Signature(obj3, private)
	if err != nil {
		return "", fmt.Errorf("sign.EIP712Signature err: %s", err.Error())
	}

	signData = append(signature, mmHash...)

	hexChainId := fmt.Sprintf("%x", chainId)
	chainIdData := common.Hex2Bytes(fmt.Sprintf("%016s", hexChainId))
	signData = append(signData, chainIdData...)
	return common.Bytes2Hex(signData), nil
}
