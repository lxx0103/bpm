package cmd

import (
	"bpm/api/v1/auth"
	"bpm/api/v1/client"
	"bpm/api/v1/component"
	"bpm/api/v1/element"
	"bpm/api/v1/event"
	"bpm/api/v1/member"
	"bpm/api/v1/node"
	"bpm/api/v1/organization"
	"bpm/api/v1/position"
	"bpm/api/v1/project"
	"bpm/api/v1/template"
	"bpm/api/v1/upload"
	"bpm/core/config"
	"bpm/core/database"
	"bpm/core/log"
	"bpm/core/router"
)

func Run(args []string) {
	config.LoadConfig(args[1])
	log.ConfigLogger()
	// cache.ConfigCache()
	database.ConfigMysql()
	// event.Subscribe(user.Subscribe, auth.Subscribe, inventory.Subscribe)
	r := router.InitRouter()
	router.InitPublicRouter(r, auth.Routers)
	router.InitAuthRouter(r, organization.Routers, project.Routers, event.Routers, component.Routers, auth.AuthRouter, client.Routers, position.Routers, member.Routers, template.Routers, node.Routers, element.Routers, upload.Routers)
	router.InitWxRouter(r, event.WxRouters, project.WxRouters, upload.WxRouters, component.WxRouters, position.WxRouters, auth.WxRouters, client.WxRouters, member.WxRouters, template.WxRouters)
	router.RunServer(r)
}
