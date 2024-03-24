package bootstrap

import (
	"crypto-follower/core/config"
	"crypto-follower/core/engine"
	"crypto-follower/core/log"
	"github.com/gookit/color"
)

var Global *globalContent

type globalContent struct {
	ApplicationConf config.ApplicationConfig
	MainEngine      *engine.MainEngine
}

func Init() {
	Global = &globalContent{}
	//初始化配置文件
	initApplicationConfig()
	//初始化日志
	initLog(&Global.ApplicationConf)
	//初始化主引擎
	initMainEngine(&Global.ApplicationConf)
	color.Greenln("启动完成。")
}

func initMainEngine(app *config.ApplicationConfig) {
	mainEngine := engine.NewMainEngine(engine.NewEventEngine())
	mainEngine.InitEngines(app)
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
