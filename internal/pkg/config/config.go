package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

var cfg Config

// Config 系统配置数据类型
type Config struct {
	Database struct {
		User     string `toml:"user"`
		Password string `toml:"password"`
		IP       string `toml:"ip"`
		Port     string `toml:"port"`
		Name     string `toml:"name"`
		MaxConn  int    `toml:"max_conn"`
	} `toml:"database"`
}

// Get 获得系统配置
func Get() Config {
	return cfg
}

// Load 载入系统配置
func Load(data []byte) error {
	return toml.Unmarshal(data, &cfg)
}

// LoadFile 从文件载入配置信息
func LoadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read file, %w", err)
	}

	return Load(data)
}
