package service

import (
	"blob-index/config"
	"blob-index/dao"
)

type Service struct {
	c *config.Config
	d *dao.Dao
}

func New(_c *config.Config, _d *dao.Dao) (service *Service) {
	service = &Service{
		c: _c,
		d: _d,
	}
	go service.loopGetTotalCount()

	return service
}
