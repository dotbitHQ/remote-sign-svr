package http_svr

import (
	"context"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"remote-sign-svr/http_svr/handle"
)

var (
	log = logger.NewLogger("http_svr", logger.LevelDebug)
)

type HttpSvr struct {
	Ctx    context.Context
	Addr   string
	H      *handle.HttpHandle
	engine *gin.Engine
	srv    *http.Server
}

func (h *HttpSvr) Run() {
	h.engine = gin.New()
	h.initRouter()
	h.srv = &http.Server{
		Addr:    h.Addr,
		Handler: h.engine,
	}
	go func() {
		if err := h.srv.ListenAndServe(); err != nil {
			log.Error("ListenAndServe err:", err)
		}
	}()
}

func (h *HttpSvr) Shutdown() {
	if h.srv != nil {
		log.Warn("HttpSvr Shutdown ... ")
		if err := h.srv.Shutdown(h.Ctx); err != nil {
			log.Error("Shutdown err:", err.Error())
		}
	}
}
