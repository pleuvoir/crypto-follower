package main

import (
	"crypto-follower/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/pleuvoir/gamine"
	"github.com/pleuvoir/gamine/component/log"
	"github.com/pleuvoir/gamine/component/restful"
	"github.com/pleuvoir/gamine/component/socketio"
	"github.com/pleuvoir/gamine/component/sqlite"
	"github.com/pleuvoir/gamine/core"
	"github.com/pleuvoir/gamine/helper/helper_config"
	"github.com/pleuvoir/gamine/helper/helper_os"
	"net/http"
	"path/filepath"
	"time"
)

var (
	profile string
	workDir string
	conf    Config
)

type Config struct {
	HttpConf   HttpConf   `yaml:"http"`
	SocketConf SocketConf `yaml:"socket"`
}

type HttpConf struct {
	Port string `yaml:"port"`
}

type SocketConf struct {
	Port int `yaml:"port"`
}

func init() {
	if profile == "" {
		profile = core.Dev
	}
	executePath, _ := helper_os.CurrentExecutePath()
	abs := helper_os.Abs(".")

	color.Greenln(fmt.Sprintf("当前可执行文件路径：%s，当前目录绝对路径：%s", executePath, abs))

	if profile == core.Prod {
		workDir = executePath
	}
	if profile == core.Dev {
		workDir = abs
	}
	conf = Config{}

	appConfigManager := core.NewConfigManager()
	appConfigManager.SetConfigName("app")
	appConfigManager.SetConfigType("yml")
	appConfigManager.AddConfigPath(filepath.Join(workDir, "configs"))
	if err := appConfigManager.LoadConfigFile(); err != nil {
		panic(err)
	}
	config := appConfigManager.GetConfig()
	if err := helper_config.InjectAnotherStructByYaml(config, &conf); err != nil {
		panic(err)
	}
}

func main() {
	gamine.SetEnvName(profile)
	gamine.SetWorkDir(workDir)
	gamine.InstallComponents(&log.Instance{}, &sqlite.Instance{})

	bindSocketIO(conf.SocketConf)
	runRestful(conf.HttpConf)
}

func bindSocketIO(socketConf SocketConf) {
	goSocketIO := socketio.New(socketConf.Port)
	_ = goSocketIO.WithRequest(func(msg socketio.RequestMessage) socketio.ResponseMessage {
		return socketio.ResponseMessage{Data: "hello world"}
	})
	_ = goSocketIO.Run()
	go func() {
		count := 1
		for true {
			//可以随时获取
			socketio.Get().PushResponse(socketio.ResponseMessage{Data: count, MethodName: "updateCount"})
			count++
			time.Sleep(3 * time.Second)
		}
	}()
}

func runRestful(restfulConf HttpConf) {
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
