package http_svr

import (
	"encoding/json"
	"github.com/dotbitHQ/das-lib/http_api"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

func (h *HttpSvr) initRouter() {
	h.engine.Use(toolib.MiddlewareCors())
	h.engine.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
	v1 := h.engine.Group("v1")
	{
		// cache
		//longExpireTime, longDataTime := time.Second*15, time.Minute*10
		//shortExpireTime, shortDataTime, lockTime := time.Second*5, time.Minute*3, time.Minute
		//cacheHandleShort := toolib.MiddlewareCacheByRedis(h.rc.GetRedisClient(), false, shortDataTime, lockTime, shortExpireTime, respHandle)
		//cacheHandleLong := toolib.MiddlewareCacheByRedis(h.rc.GetRedisClient(), false, longDataTime, lockTime, longExpireTime, respHandle)
		//cacheHandleShortCookies := toolib.MiddlewareCacheByRedis(h.rc.GetRedisClient(), true, shortDataTime, lockTime, shortExpireTime, respHandle)

		// query
		v1.POST("/version", DoMonitorLog("version"), h.H.Version)
		v1.POST("/svr/info", DoMonitorLog("svr_info"), h.H.SvrInfo)
		v1.POST("/remote/sign", DoMonitorLog("remote_sign"), h.H.RemoteSign)
		v1.POST("/address/info", DoMonitorLog("address_info"), h.H.AddressInfo)
		v1.POST("/address/disable", DoMonitorLog("address_disable"), h.H.AddressDisable)

		// operate
		v1.POST("/init/svr", DoMonitorLog("init_svr"), h.H.InitSvr)
		v1.POST("/import/address", DoMonitorLog("import_address"), h.H.ImportAddress)
	}
}

func respHandle(c *gin.Context, res string, err error) {
	if err != nil {
		log.Error("respHandle err:", err.Error())
		c.AbortWithStatusJSON(http.StatusOK, http_api.ApiRespErr(http.StatusInternalServerError, err.Error()))
	} else if res != "" {
		var respMap map[string]interface{}
		_ = json.Unmarshal([]byte(res), &respMap)
		c.AbortWithStatusJSON(http.StatusOK, respMap)
	}
}
