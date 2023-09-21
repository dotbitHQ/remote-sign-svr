package handle

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"remote-sign-svr/config"
	"remote-sign-svr/tables"
	"strings"
)

type ReqAddressInfo struct {
	Address string `json:"address"`
}

type RespAddressInfo struct {
	AddrChain    tables.AddrChain `json:"addr_chain"`
	Address      string
	Private      string
	AddrStatus   tables.AddrStatus `json:"addr_status"`
	Remark       string
	CompressType tables.CompressType `json:"compress_type"`
}

func (h *HttpHandle) AddressInfo(ctx *gin.Context) {
	var (
		funcName             = "AddressInfo"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqAddressInfo
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

	if err = h.doAddressInfo(&req, &apiResp); err != nil {
		log.Error("doAddressInfo err:", err.Error(), funcName, clientIp, remoteAddr, ctx)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func checkIP(apiResp *http_api.ApiResp, clientIp, remoteAddr string) {
	if strings.Contains(clientIp, ":") {
		clientIp = clientIp[:strings.Index(clientIp, ":")]
	}
	if strings.Contains(remoteAddr, ":") {
		remoteAddr = remoteAddr[:strings.Index(remoteAddr, ":")]
	}
	if _, ok := config.Cfg.IpWhitelist[clientIp]; !ok {
		if _, ok = config.Cfg.IpWhitelist[remoteAddr]; !ok {
			apiResp.ApiRespErr(http_api.ApiCodeIpBlockingAccess, "IP Blocking Access")
		}
	}
}

func (h *HttpHandle) doAddressInfo(req *ReqAddressInfo, apiResp *http_api.ApiResp) error {
	var resp RespAddressInfo

	// Check if the service is active
	if key := config.Cfg.GetKey(); key == "" {
		apiResp.ApiRespErr(http_api.ApiCodeServiceNotActivated, "service not activated")
		return nil
	}

	addrInfo, err := h.DbDao.GetAddressInfo(req.Address)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, err.Error())
		return fmt.Errorf("GetAddressInfo err: %s", err.Error())
	} else if addrInfo.Id == 0 {
		apiResp.ApiRespErr(http_api.ApiCodeWalletAddrNotExist, "Wallet address does not exist")
		return nil
	}
	resp.Address = addrInfo.Address
	resp.AddrChain = addrInfo.AddrChain
	resp.AddrStatus = addrInfo.AddrStatus
	resp.Private = addrInfo.Private
	resp.Remark = addrInfo.Remark
	resp.CompressType = addrInfo.CompressType

	apiResp.ApiRespOK(resp)
	return nil
}
