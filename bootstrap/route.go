package bootstrap

import (
	"GoBlog/pkg/route"
	"GoBlog/routes"

	"github.com/gorilla/mux"
)

//SetupRoute 路由初始化
func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)

	route.SetRoute(router)

	return router
}