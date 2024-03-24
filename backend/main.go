package main

import (
	"crypto-follower/bootstrap"
	"crypto-follower/core/config"
	"crypto-follower/core/helper"
	"github.com/gookit/color"
	"os"
	"path/filepath"
)

func init() {
	rootPath, err := helper.RootPath()
	if err != nil {
		panic(err)
	}
	confPath := filepath.Join(rootPath, "/bin/app.yml")
	err = os.Setenv(config.ApplicationEnvVar, confPath)
	if err != nil {
		panic(err)
	}
	color.Greenf("[本地使用]设置配置文件环境变量，%s", confPath)
	println()
}

func main() {
	bootstrap.Init()
}
