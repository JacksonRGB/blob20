package dao

import (
	"blob-index/config"
	dbmodel "blob-index/model/db"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Dao struct {
	c         *config.Config
	ethClient *ethclient.Client
	ethRPC    *rpc.Client
	db        *gorm.DB
}

func New(_c *config.Config) (dao *Dao, err error) {
	dao = &Dao{
		c: _c,
	}
	dao.ethClient, err = ethclient.Dial(_c.Chain.RPC)
	if err != nil {
		return
	}
	dao.ethRPC, err = rpc.Dial(_c.Chain.RPC)
	if err != nil {
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		_c.Mysql.User, _c.Mysql.Password, _c.Mysql.Host, _c.Mysql.Port, _c.Mysql.Database)
	dao.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return
	}
	sqlDB, err := dao.db.DB()
	if err != nil {
		return
	}
	sqlDB.SetMaxOpenConns(_c.Mysql.MaxConn)
	sqlDB.SetMaxIdleConns(_c.Mysql.MaxIdleConn)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	if _c.Mysql.Migrate {
		err = dao.db.AutoMigrate(&dbmodel.Mint{}, &dbmodel.Height{})
		if err != nil {
			return
		}
	}

	return dao, nil
}
