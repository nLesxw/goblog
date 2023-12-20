package controllers

import (
	"GoBlog/pkg/logger"
	"fmt"
	"html/template"
	"net/http"
)

//PagesController 处理静态页面
type PagesController struct{

}

//Home 首页
func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
    
    fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

//About 关于我们页面
func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("views/about.gohtml")
    logger.LogError(err)

    err = tmpl.Execute(w, "n1_esxw")
    logger.LogError(err)
}

//NotFound 404页面
func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
    
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}