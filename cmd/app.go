package cmd

import (
	"bpm/api/v1/auth"
	"bpm/api/v1/component"
	"bpm/api/v1/event"
	"bpm/api/v1/organization"
	"bpm/api/v1/project"
	"bpm/core/cache"
	"bpm/core/config"
	"bpm/core/database"
	"bpm/core/log"
	"bpm/core/router"
)

func Run() {
	config.LoadConfig("config.toml")
	log.ConfigLogger()
	cache.ConfigCache()
	database.ConfigMysql()
	// event.Subscribe(user.Subscribe, auth.Subscribe, inventory.Subscribe)
	r := router.InitRouter()
	router.InitPublicRouter(r, auth.Routers)
	router.InitAuthRouter(r, organization.Routers, project.Routers, event.Routers, component.Routers, auth.AuthRouter)
	router.RunServer(r)
}
