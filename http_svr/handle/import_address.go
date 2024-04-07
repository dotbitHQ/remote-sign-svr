package handle

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"regexp"
	"remote-sign-svr/config"
	"remote-sign-svr/encrypt"
	"remote-sign-svr/tables"
)

type ReqImportAddress struct {
	AddrChain tables.AddrChain `json:"addr_chain"`
	Address   string           `json:"address"`
	Private   string           `json:"private"`
	Remark    string           `json:"remark"`
}

type RespImportAddress struct {
}

func (h *HttpHandle) ImportAddress(ctx *gin.Context) {
	var (
		funcName             = "ImportAddress"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqImportAddress
		apiResp              http_api.ApiResp
		err                  error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, clientIp, remoteAddr, ctx)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, remoteAddr, toolib.JsonString(req), ctx)

	checkIP(&apiResp, clientIp, remoteAddr)
	if apiResp.ErrNo != http_api.ApiCodeSuccess {
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	if err = h.doImportAddress(&req, &apiResp); err != nil {
		log.Error("doImportAddress err:", err.Error(), funcName, clientIp, remoteAddr, ctx)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doImportAddress(req *ReqImportAddress, apiResp *http_api.ApiResp) error {
	var resp RespImportAddress

	// Check if the service is active
	if key := config.Cfg.GetKey(); key == "" {
		apiResp.ApiRespErr(http_api.ApiCodeServiceNotActivated, "service not activated")
		return nil
	}

	// Check that the private key is properly encrypted
	private, err := encrypt.AesDecrypt(req.Private, config.Cfg.GetKey())
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeKeyDiff, "Invalid encryption private key")
		return fmt.Errorf("encrypt.AesDecrypt err: %s", err.Error())
	}

	compressType := tables.CompressTypeFalse
	switch req.AddrChain {
	case tables.AddrChainEVM:
		if ok, err := regexp.MatchString("^0x[0-9a-fA-F]{40}$", req.Address); err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, err.Error())
			return fmt.Errorf("regexp.MatchString err: %s", err.Error())
		} else if !ok {
			apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address invalid")
			return nil
		}
	case tables.AddrChainTRON:
		if _, err := common.TronBase58ToHex(req.Address); err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, err.Error())
			return fmt.Errorf("common.TronBase58ToHex err: %s", err.Error())
		}
	case tables.AddrChainDOGE:
		if _, err := common.Base58CheckDecode(req.Address, common.DogeCoinBase58Version); err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, err.Error())
			return fmt.Errorf("common.Base58CheckDecode err: %s", err.Error())
		}
		_, _, compress, err := bitcoin.HexPrivateKeyToScript(req.Address, bitcoin.GetDogeMainNetParams(), private)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, err.Error())
			return fmt.Errorf("bitcoin.HexPrivateKeyToScript err: %s", err.Error())
		}
		if compress {
			compressType = tables.CompressTypeTrue
		}
	case tables.AddrChainBTC:
		netParams, _, _, err := bitcoin.FormatBTCAddr(req.Address)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, err.Error())
			return fmt.Errorf("bitcoin.FormatBTCAddr err: %s", err.Error())
		}
		_, _, compress, err := bitcoin.HexPrivateKeyToScript(req.Address, netParams, private)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, err.Error())
			return fmt.Errorf("bitcoin.HexPrivateKeyToScript err: %s", err.Error())
		}
		if compress {
			compressType = tables.CompressTypeTrue
		}
	case tables.AddrChainCKB:
		if _, err := address.Parse(req.Address); err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, err.Error())
			return fmt.Errorf("address.Parse err: %s", err.Error())
		}
	default:
		apiResp.ApiRespErr(http_api.ApiCodeUnsupportedAddrChain, "unsupported address chain")
		return nil
	}

	addrInfo := tables.TableAddressInfo{
		Id:           0,
		AddrChain:    req.AddrChain,
		Address:      req.Address,
		Private:      req.Private,
		AddrStatus:   tables.AddrStatusDefault,
		Remark:       req.Remark,
		CompressType: compressType,
	}
	if err := h.DbDao.CreateAddressInfo(addrInfo); err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, err.Error())
		return fmt.Errorf("CreateAddressInfo err: %s", err.Error())
	}

	apiResp.ApiRespOK(resp)
	return nil
}
