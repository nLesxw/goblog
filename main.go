package main

import (
	"GoBlog/bootstrap"
	"GoBlog/pkg/database"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var router *mux.Router



func forceHTMLMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //1.设置标头
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        //2.继续处理
        next.ServeHTTP(w, r)
    })
}

func removeTrailingSlash(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //1.除首页外，移除所有路径后带"/"
        if r.URL.Path != "/" {
            r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
            
        }
        //2.将请求传递下去
        next.ServeHTTP(w, r)
    })
}

func main() {
    database.Initialize()
    _ = database.DB

    bootstrap.SetupDB()
    router = bootstrap.SetupRoute()

    

    //中间件：强制内容类型为 HTML
    router.Use(forceHTMLMiddleware)

    // 通过命名路由获取 URL 示例
    
    http.ListenAndServe(":3000", removeTrailingSlash(router))
}