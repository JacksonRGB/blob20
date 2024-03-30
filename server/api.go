package server

import (
	"blob-index/config"
	"blob-index/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var srv *service.Service
var conf *config.Config

func Run(_srv *service.Service, _conf *config.Config) {
	srv = _srv
	conf = _conf
	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())

	_cors := cors.DefaultConfig()
	_cors.AllowAllOrigins = true
	engine.Use(cors.New(_cors))

	router(engine)
	log.Infof("start http server listening %s", conf.Server.Listen)
	if err := engine.Run(conf.Server.Listen); err != nil {
		log.Error("http server run error: ", err)
	}
}
