package main

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/jinzhu/gorm"
	"github.com/securingsincity/gizmo-boilerplate/service"
)

func main() {
	var cfg *service.Config
	config.LoadJSONFile("./config.json", &cfg)
	server.Init("gizmo-boilerplate", cfg.Server)
	dbSql, err := cfg.MySQL.DB()
	if err != nil {
		server.Log.Fatal("unable to connect to mysql ", err)
	}
	cfg.DB, _ = gorm.Open("mysql", dbSql)
	err = server.Register(service.NewJSONService(cfg))
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
