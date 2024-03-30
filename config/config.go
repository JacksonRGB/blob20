package config

import (
	"flag"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Debug  bool         `toml:"debug"`
	Chain  ChainConfig  `toml:"chain"`
	Mysql  MysqlConfig  `toml:"mysql"`
	Server ServerConfig `toml:"server"`
}

type ChainConfig struct {
	RPC string `toml:"rpc"`
	ID  int    `toml:"id"`
}

type MysqlConfig struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	User        string `toml:"user"`
	Password    string `toml:"password"`
	Database    string `toml:"database"`
	MaxConn     int    `toml:"max_conn"`
	MaxIdleConn int    `toml:"max_idle_conn"`
	Migrate     bool   `toml:"migrate"`
}

type ServerConfig struct {
	Listen  string `toml:"listen"`
	Disable bool   `toml:"disable"`
}

var confPath = flag.String("c", "config.toml", "config file path")

func New() (config *Config, err error) {
	config = new(Config)
	_, err = toml.DecodeFile(*confPath, config)
	return
}
