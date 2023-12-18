package main

import (
	"GoBlog/app/middlewares"
	"GoBlog/bootstrap"
	"GoBlog/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router

func main() {
    bootstrap.SetupDB()
    router = bootstrap.SetupRoute()

    err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
    logger.LogError(err)
}