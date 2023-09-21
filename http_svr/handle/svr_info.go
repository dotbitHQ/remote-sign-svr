package handle

import (
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"remote-sign-svr/config"
)

type ReqSvrInfo struct {
}

type RespSvrInfo struct {
	Activated bool `json:"activated"`
}

func (h *HttpHandle) SvrInfo(ctx *gin.Context) {
	var (
		funcName             = "SvrInfo"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqSvrInfo
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

	if err = h.doSvrInfo(&req, &apiResp); err != nil {
		log.Error("doSvrInfo err:", err.Error(), funcName, clientIp, remoteAddr, ctx)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doSvrInfo(req *ReqSvrInfo, apiResp *http_api.ApiResp) error {
	var resp RespSvrInfo

	if key := config.Cfg.GetKey(); key != "" {
		resp.Activated = true
	}

	apiResp.ApiRespOK(resp)
	return nil
}
