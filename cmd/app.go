package cmd

import (
	"bpm/api/v1/assignment"
	"bpm/api/v1/auth"
	"bpm/api/v1/client"
	"bpm/api/v1/common"
	"bpm/api/v1/component"
	"bpm/api/v1/element"
	"bpm/api/v1/event"
	"bpm/api/v1/example"
	"bpm/api/v1/meeting"
	"bpm/api/v1/member"
	"bpm/api/v1/message"
	"bpm/api/v1/node"
	"bpm/api/v1/organization"
	"bpm/api/v1/position"
	"bpm/api/v1/project"
	"bpm/api/v1/template"
	"bpm/api/v1/upload"
	"bpm/api/v1/vendors"
	"bpm/core/config"
	"bpm/core/database"
	event2 "bpm/core/event"
	"bpm/core/log"
	"bpm/core/router"
)

func Run(args []string) {
	config.LoadConfig(args[1])
	log.ConfigLogger()
	// cache.ConfigCache()
	database.ConfigMysql()
	event2.Subscribe(message.Subscribe, event.Subscribe)
	r := router.InitRouter()
	router.InitPublicRouter(r, auth.Routers, organization.PortalRouters, example.PortalRouters, vendors.PortalRouters, common.PortalRouters, project.PortalRouters)
	router.InitAuthRouter(r, organization.Routers, project.Routers, event.Routers, component.Routers, auth.AuthRouter, client.Routers, position.Routers, member.Routers, template.Routers, node.Routers, element.Routers, upload.Routers, example.Routers, common.Routers, vendors.Routers, meeting.Routers, assignment.Routers)
	router.InitWxRouter(r, event.WxRouters, project.WxRouters, upload.WxRouters, component.WxRouters, position.WxRouters, auth.WxRouters, client.WxRouters, member.WxRouters, template.WxRouters, example.WxRouters, organization.WxRouters, meeting.WxRouters, assignment.WxRouters)
	router.RunServer(r)
}
