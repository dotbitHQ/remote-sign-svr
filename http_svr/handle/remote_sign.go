package handle

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"math/big"
	"net/http"
	"remote-sign-svr/config"
	"remote-sign-svr/encrypt"
	"remote-sign-svr/tables"
	"remote-sign-svr/wallet"
)

type ReqRemoteSign struct {
	SignType   wallet.SignType   `json:"sign_type"`
	Address    string            `json:"address"`
	EvmChainID int64             `json:"evm_chain_id"`
	Data       string            `json:"data"`
	MMJson     *common.MMJsonObj `json:"mm_json"`
}

func (r *ReqRemoteSign) ChainId() *big.Int {
	return big.NewInt(r.EvmChainID)
}

type RespRemoteSign struct {
	Data string `json:"data"`
}

func (h *HttpHandle) RemoteSign(ctx *gin.Context) {
	var (
		funcName             = "RemoteSign"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqRemoteSign
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

	if err = h.doRemoteSign(&req, &apiResp); err != nil {
		log.Error("doRemoteSign err:", err.Error(), funcName, clientIp, remoteAddr, ctx)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doRemoteSign(req *ReqRemoteSign, apiResp *http_api.ApiResp) error {
	var resp RespRemoteSign

	// Check if the service is active
	key := config.Cfg.GetKey()
	if key == "" {
		apiResp.ApiRespErr(http_api.ApiCodeServiceNotActivated, "service not activated")
		return nil
	}

	// Get address info by addr
	addrInfo, err := h.DbDao.GetAddressInfo(req.Address)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, err.Error())
		return fmt.Errorf("GetAddressInfo err: %s", err.Error())
	} else if addrInfo.Id == 0 {
		apiResp.ApiRespErr(http_api.ApiCodeWalletAddrNotExist, "Wallet address does not exist")
		return nil
	}
	if addrInfo.AddrStatus != tables.AddrStatusDefault {
		apiResp.ApiRespErr(http_api.ApiCodeAddressStatusNotNormal, fmt.Sprintf("address status: %d", addrInfo.AddrStatus))
		return nil
	}

	//
	addrInfo.Private, err = encrypt.AesDecrypt(addrInfo.Private, key)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, err.Error())
		return fmt.Errorf("encrypt.AesDecrypt err: %s", err.Error())
	}

	// sign
	data, err := wallet.Sign(req.SignType, addrInfo, req.Data, req.ChainId(), req.MMJson)
	if err != nil {
		switch err {
		case wallet.ErrUnsupportedAddrChain:
			apiResp.ApiRespErr(http_api.ApiCodeUnsupportedAddrChain, fmt.Sprintf("wallet.Sign err: %s", err.Error()))
			return nil
		case wallet.ErrUnsupportedSignType:
			apiResp.ApiRespErr(http_api.ApiCodeUnsupportedAddrChain, fmt.Sprintf("wallet.Sign err: %s", err.Error()))
			return nil
		default:
			apiResp.ApiRespErr(http_api.ApiCodeError500, fmt.Sprintf("wallet.Sign err: %s", err.Error()))
			return fmt.Errorf("wallet.Sign err: %s", err.Error())
		}
	}

	resp.Data = data
	apiResp.ApiRespOK(resp)
	return nil
}
