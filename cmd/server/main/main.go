package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/HuangXiaoL/xiaoshuo/internal/pkg/novelread"

	"github.com/HuangXiaoL/xiaoshuo/internal/pkg/config"
	"github.com/HuangXiaoL/xiaoshuo/internal/pkg/connection"
	"github.com/sirupsen/logrus"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "", "config file")
	flag.Parse()

	initLog()
	if err := initConfig(); err != nil {
		logrus.WithError(err).Fatal("load config")
	}

	if err := connection.Init(); err != nil {
		logrus.WithError(err).Fatal("initialize connection package")
	}
}

func initLog() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if s := os.Getenv("LOG_LEVEL"); s != "" {
		if lvl, err := logrus.ParseLevel(s); err == nil {
			logrus.SetLevel(lvl)
		}
	}
}

//指定配置文件
func initConfig() error {
	if configFile == "" {
		return fmt.Errorf("require file")
	}
	return config.LoadFile(configFile)
}
func main() {
	st := time.Now()
	novelread.NovelRead()
	useTime := time.Since(st)
	logrus.Printf("用时为：%s", useTime)
}
