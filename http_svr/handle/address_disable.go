package handle

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"remote-sign-svr/config"
	"remote-sign-svr/tables"
)

type ReqAddressDisable struct {
	Address    string            `json:"address"`
	AddrStatus tables.AddrStatus `json:"addr_status"`
}

type RespAddressDisable struct {
}

func (h *HttpHandle) AddressDisable(ctx *gin.Context) {
	var (
		funcName             = "AddressDisable"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqAddressDisable
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

	if err = h.doAddressDisable(&req, &apiResp); err != nil {
		log.Error("doAddressDisable err:", err.Error(), funcName, clientIp, remoteAddr, ctx)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAddressDisable(req *ReqAddressDisable, apiResp *http_api.ApiResp) error {
	var resp RespAddressDisable

	// Check if the service is active
	if key := config.Cfg.GetKey(); key == "" {
		apiResp.ApiRespErr(http_api.ApiCodeServiceNotActivated, "service not activated")
		return nil
	}

	if err := h.DbDao.UpdateAddressStatus(req.Address, req.AddrStatus); err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, err.Error())
		return fmt.Errorf("UpdateAddressStatus err: %s", err.Error())
	}

	apiResp.ApiRespOK(resp)
	return nil
}
