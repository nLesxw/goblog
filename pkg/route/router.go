package route

import (
	"GoBlog/pkg/config"
	"GoBlog/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

var route *mux.Router

//SetRoute 设置路由实例，以供 Name2URL 等函数使用
func SetRoute(r *mux.Router) {
    route = r
}

//RouteName2URL 通过路由名称来获取URL
func Name2URL(routeName string, pairs ...string) string{
    url, err := route.Get(routeName).URL(pairs...)
    if err != nil {
        logger.LogError(err)
        return ""
    }

    return config.GetString("app.url") + url.String()
}

func GetRouteVariable(parameterName string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[parameterName]
}