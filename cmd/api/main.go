package main

import (
	"blob-index/config"
	"blob-index/dao"
	"blob-index/server"
	"blob-index/service"
	"flag"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

var migrate = flag.Bool("migrate", false, "migrate database")

func main() {
	flag.Parse()
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	conf.Mysql.Migrate = *migrate

	da, err := dao.New(conf)
	if err != nil {
		panic(err)
	}

	if conf.Debug {
		log.SetLevel(log.DebugLevel)
	}

	svs := service.New(conf, da)

	server.Run(svs, conf)
}
