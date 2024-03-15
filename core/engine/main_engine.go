package engine

import (
	"crypto-follower/core/event"
	"github.com/gookit/color"
	"sync"
	"time"
)

type Engineer interface {
	Name() string
	Start()
	Stop()
}

const DefaultMainEngineName = "main-engine"

type MainEngine struct {
	TodayDate   string
	eventEngine *event.Engine
	engineMap   sync.Map //[string]Engineer 引擎合集
}

// NewMainEngine 构建主引擎
func NewMainEngine(eventEngine *event.Engine) *MainEngine {
	mainEngine := MainEngine{}
	mainEngine.TodayDate = time.Now().Format("2006-01-02")
	mainEngine.eventEngine = eventEngine
	mainEngine.engineMap = sync.Map{}
	return &mainEngine
}

func (m *MainEngine) InitEngines() {
}

// RegisterListener 注册事件
func (m *MainEngine) RegisterListener(t event.Type, f func(e event.Event)) {
	m.eventEngine.Register(t, event.AdaptEventHandlerFunc(f))
}

func (m *MainEngine) Name() string {
	return DefaultMainEngineName
}

// Start 主引擎启动
func (m *MainEngine) Start() {
	//启动事件引擎
	m.eventEngine.StartAll()
	for _, engine := range m.GetAllEngine() {
		engine.Start()
	}
	color.Greenln("主引擎已启动")
}

// Stop 主引擎关闭
func (m *MainEngine) Stop() {
	//关闭事件引擎
	m.eventEngine.StopAll()
	//关闭所有引擎
	for _, engine := range m.GetAllEngine() {
		engine.Stop()
	}
	color.Redln("主引擎已关闭")
}

// AddEngine 增加引擎
func (m *MainEngine) AddEngine(engine Engineer) {
	m.engineMap.Store(engine.Name(), engine)
}

// GetEngine 获取引擎
func (m *MainEngine) GetEngine(engineName string) Engineer {
	e, ok := m.engineMap.Load(engineName)
	if ok {
		engine := e.(Engineer)
		return engine
	}
	return nil
}

func (m *MainEngine) GetAllEngine() (engines []Engineer) {
	r := make(map[string]Engineer)
	m.engineMap.Range(func(key, value any) bool {
		r[key.(string)] = value.(Engineer)
		return true
	})
	for _, engine := range r {
		engines = append(engines, engine)
	}
	return engines
}
