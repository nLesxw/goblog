package main

import (
	"GoBlog/app/middlewares"
	"GoBlog/bootstrap"
	"GoBlog/pkg/logger"
	"net/http"
	c "GoBlog/pkg/config"

	"github.com/gorilla/mux"
)

var router *mux.Router

func init() {
	c.Initialize()
}

func main() {
	// 初始化 SQL
    bootstrap.SetupDB()

	// 初始化路由绑定
    router = bootstrap.SetupRoute()

    err := http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
    logger.LogError(err)
}