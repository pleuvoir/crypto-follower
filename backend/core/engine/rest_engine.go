package engine

import (
	"context"
	"crypto-follower/core/helper"
	"crypto-follower/core/log"
	"crypto-follower/restful/controller"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const RestEngineName = "gin"

// ServerState 服务状态
type ServerState int

const (
	// ServerStarting 服务启动中
	ServerStarting ServerState = 1
	// ServerStarted 服务已启动
	ServerStarted ServerState = 2
	// ServerFailed 服务启动失败
	ServerFailed ServerState = 3
)

// ServerStartedListener 服务启动后调用的监听程序
type ServerStartedListener func(engine *RestEngine)

type RestEngine struct {
	gin         *gin.Engine
	httpServer  *http.Server
	startedChan chan bool
	//延时测试端口的时间
	testPortDelayed time.Duration
	//测试端口的重试次数，若设置为小于1的数则按1次处理
	testPortRetryTimes int
	//服务启动后调用的监听程序
	serverStartedListener ServerStartedListener
	// 端口号
	port string
	// 服务状态
	State ServerState
}

func NewRestEngine(port string, listener ServerStartedListener) *RestEngine {
	return &RestEngine{
		startedChan:           make(chan bool, 1),
		testPortDelayed:       time.Second * 2,
		testPortRetryTimes:    3,
		port:                  port,
		serverStartedListener: listener,
	}
}

func (e *RestEngine) Name() string {
	return RestEngineName
}

func (e *RestEngine) Start() {
	gin.SetMode(gin.ReleaseMode)
	e.gin = gin.New()
	engine := e.gin
	engine.Use(gin.Recovery())
	engine.Use(Logger())
	Validator()
	Router(engine)
	e.State = ServerStarting
	serv := &http.Server{Addr: ":" + e.port, Handler: engine}
	e.httpServer = serv
	e.startServer(serv) //异步启动
	e.listenServerStarted()
	e.gracefulShutdown(serv) //阻塞进程等待退出

}

func (e *RestEngine) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := e.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server exiting")
}

func Router(engine *gin.Engine) {
	executePath, err := currentExecutePath()
	if err != nil {
		panic(err)
	}
	frontendPath := filepath.Join(executePath, "../frontend")
	if !helper.IsExists(frontendPath) {
		color.Yellowf("gin静态文件打包路径，可执行路径下未找到 %s", frontendPath)
		color.Yellowln()
		rootPath, _ := helper.RootPath()
		frontendPath = filepath.Join(rootPath, "frontend")
		color.Yellowf("尝试使用项目根路径..%s", frontendPath)
		color.Yellowln()
	}

	if !helper.IsExists(frontendPath) {
		color.Redln("gin启动失败，未找到静态资源文件路径。")
		panic("gin启动失败，未找到静态资源文件路径")
	}

	color.Greenf("gin静态文件打包路径 %s", frontendPath)
	color.Greenln()
	engine.StaticFS("app", http.Dir(frontendPath))
	engine.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(frontendPath, "index.html"))
	})
	index := engine.Group("/")
	{
		indexController := controller.NewIndexController()
		index.GET("/welcome", indexController.Welcome)
		index.GET("/", indexController.Welcome)
	}
}

func currentExecutePath() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(dir), nil
}

func Validator() {

}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		log.Infof("| %3d | %13v | %15s | %s | %s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func (e *RestEngine) startServer(serv *http.Server) {
	go func() {
		e.startedChan <- true
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.State = ServerFailed
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

// listenServerStarted 启动监听，在服务启动时检测http服务监听的端口
func (e *RestEngine) listenServerStarted() {
	if e.serverStartedListener == nil {
		return
	}
	go func() {
		<-e.startedChan
		if e.testPortRetry() {
			e.State = ServerStarted
			log.Info("Server started.")
			e.serverStartedListener(e)
		}
	}()
}

// testPortRetry 检测http服务监听的端口，该方法会延时阻塞执行，若监测端口超时则会重试
func (e *RestEngine) testPortRetry() bool {
	time.Sleep(e.testPortDelayed)
	testPortTimes := helper.If(e.testPortRetryTimes < 1, 1, e.testPortRetryTimes).(int)
	for i := 0; i < testPortTimes; i++ {
		if e.testPort() {
			return true
		}
	}
	return false
}

// testPort 检测一次http服务监听的端口
func (e *RestEngine) testPort() bool {
	conn, err := net.DialTimeout("tcp", ":"+e.port, time.Millisecond*500)
	defer helper.CloseQuietly(conn)
	return err == nil
}

// gracefulShutdown 优雅关闭服务
func (e *RestEngine) gracefulShutdown(serv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := serv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server exiting")
}
