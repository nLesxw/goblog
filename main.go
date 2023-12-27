package main

import (
	"GoBlog/app/middlewares"
	"GoBlog/bootstrap"
	"GoBlog/config"
	c "GoBlog/pkg/config"
	"GoBlog/pkg/logger"
	"net/http"
)

func init() {
	config.Initialize()
}

func main() {
	// 初始化 SQL
    bootstrap.SetupDB()

	// 初始化路由绑定
    router := bootstrap.SetupRoute()

    err := http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
    logger.LogError(err)
}