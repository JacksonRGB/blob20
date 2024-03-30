package main

import (
	"blob-index/config"
	"blob-index/dao"
	"blob-index/sync"
	"flag"
	"time"

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

	sy := sync.New(conf, da)
	for {
		sy.Start()
		log.Error("sync error, retry after 20 seconds")
		time.Sleep(time.Second * 20)
	}
}
