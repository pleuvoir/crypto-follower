package bootstrap

import (
	"crypto-follower/core/config"
	"crypto-follower/core/engine"
	"crypto-follower/core/event"
	"crypto-follower/core/log"
	"crypto-follower/restful/controller"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
)

var Global *globalContent

type globalContent struct {
	ApplicationConf config.ApplicationConfig
	MainEngine      *engine.MainEngine
	RestfulEngine   *gin.Engine
}

func Init() {
	Global = &globalContent{}
	//初始化配置文件
	initApplicationConfig()
	//初始化日志
	initLog(&Global.ApplicationConf)
	//初始化主引擎
	initMainEngine(&Global.ApplicationConf)
	//初始化restful
	initRestfulEngine(&Global.ApplicationConf)
	color.Greenln("启动完成。")
}

func initMainEngine(app *config.ApplicationConfig) {
	mainEngine := engine.NewMainEngine(event.NewEventEngine())
	mainEngine.InitEngines()
	mainEngine.Start()
	Global.MainEngine = mainEngine
}

func initApplicationConfig() {
	conf := config.NewApplicationConf()
	if err := conf.Load(); err != nil {
		panic(err)
	}
	Global.ApplicationConf = conf
}

func initLog(app *config.ApplicationConfig) {
	log.InitLog(app)
}

func initRestfulEngine(app *config.ApplicationConfig) {
	controller.SetMode(app)
	Global.RestfulEngine = gin.New()
	e := Global.RestfulEngine
	e.Use(gin.Recovery())
	e.Use(controller.Logger())
	controller.Validator()
	controller.Router(e)
	color.Greenf("Restful服务已启动 listen on http://127.0.0.1:%s", app.Server.Port)
	color.Println()
}
