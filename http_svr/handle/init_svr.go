package handle

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"net/http"
	"remote-sign-svr/config"
	"remote-sign-svr/encrypt"
)

type ReqInitSvr struct {
	Key string `json:"key"`
}

type RespInitSvr struct {
}

func (h *HttpHandle) InitSvr(ctx *gin.Context) {
	var (
		funcName             = "InitSvr"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqInitSvr
		apiResp              http_api.ApiResp
		err                  error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, clientIp, remoteAddr, ctx)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, remoteAddr, ctx)

	if err = h.doInitSvr(&req, &apiResp); err != nil {
		log.Error("doInitSvr err:", err.Error(), funcName, clientIp, remoteAddr, ctx)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doInitSvr(req *ReqInitSvr, apiResp *http_api.ApiResp) error {
	var resp RespInitSvr

	// Check if the key is correct
	list, err := h.DbDao.GetAddressListGroupByAddrChain()
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "Failed to get address list")
		return fmt.Errorf("GetAddressListGroupByAddrChain err: %s", err.Error())
	}
	for _, v := range list {
		_, err = encrypt.AesDecrypt(v.Private, req.Key)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeKeyDiff, "The current key is not the same as the original key")
			return fmt.Errorf("encrypt.AesDecrypt err: %s", err.Error())
		}
	}
	config.Cfg.SetKey(req.Key)

	apiResp.ApiRespOK(resp)
	return nil
}
