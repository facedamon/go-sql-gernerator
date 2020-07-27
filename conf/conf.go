package conf

import (
	"github.com/BurntSushi/toml"
	"path/filepath"
	"sync"
)

type db struct {
	Ip      string `toml:"ip"`
	Port    int    `toml:"port"`
	Schema  string `toml:"schema"`
	User    string `toml:"user"`
	Pwd     string `toml:"pwd"`
	MaxConn int    `toml:"maxConn"`
	MaxIdle int    `toml:"maxIdle"`
	Enable  bool   `toml:"enable"`
}

type tomlConf struct {
	Db db `toml:"db"`
}

var (
	cfg     *tomlConf
	once    sync.Once
	cfgLock sync.RWMutex
)

func Config() *tomlConf {
	once.Do(ReloadConfig)
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return cfg
}

func ReloadConfig() {
	cfgLock.Lock()
	defer cfgLock.Unlock()

	filePath, err := filepath.Abs("./conf.toml")
	if err != nil {
		panic(err)
	}
	config := new(tomlConf)
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		panic(err)
	}
	cfg = config
}
