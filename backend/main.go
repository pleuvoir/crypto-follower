package main

import (
	"crypto-follower/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/pleuvoir/gamine"
	"github.com/pleuvoir/gamine/component/log"
	"github.com/pleuvoir/gamine/component/restful"
	"github.com/pleuvoir/gamine/component/sqlite"
	"github.com/pleuvoir/gamine/core"
	"github.com/pleuvoir/gamine/helper/helper_config"
	"github.com/pleuvoir/gamine/helper/helper_os"
	"net/http"
	"path/filepath"
)

var (
	profile string
	workDir string
	conf    Config
)

type Config struct {
	RestfulConf RestfulConf `yaml:"restful"`
}

type RestfulConf struct {
	Port string `yaml:"port"`
}

func init() {
	if profile == "" {
		profile = core.Dev
	}
	if profile == core.Prod {
		cur, _ := helper_os.CurrentExecutePath()
		workDir = filepath.Join(cur, "../bin")
	}
	if profile == core.Dev {
		workDir = helper_os.Abs("../bin")
	}
	conf = Config{}
	if err := helper_config.ParseYamlStringFromPath2Struct(filepath.Join(workDir, "app.yml"), &conf); err != nil {
		panic(err)
	}
}

func main() {
	gamine.SetEnvName(profile)
	gamine.SetWorkDir(workDir)
	gamine.InstallComponents(&log.Instance{}, &sqlite.Instance{})
	runRestful(conf.RestfulConf)
}

func runRestful(restfulConf RestfulConf) {
	server := restful.NewRestServer(restfulConf.Port)
	staticPath := helper_os.Abs("../frontend")
	color.Greenln(fmt.Sprintf("静态文件目录:%s", staticPath))
	server.WithUseRequestLog(log.GetDefault())
	server.WithGinConfig(func(engine *gin.Engine) {
		engine.StaticFS("app", http.Dir(staticPath))
		engine.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(staticPath, "index.html"))
		})
	})
	server.WithGinConfig(func(e *gin.Engine) {
		index := e.Group("/")
		{
			indexController := controller.IndexController{}
			index.GET("/", indexController.Index)
		}
	})
	server.Run()
}
