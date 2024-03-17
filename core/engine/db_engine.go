package engine

import (
	"crypto-follower/core/config"
	"github.com/gookit/color"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DefaultSqlEngineName = "db-sqlite-engine"

type SqlLiteEngine struct {
	name string
	url  string
	db   *gorm.DB
}

func NewDbEngine(database *config.Database) *SqlLiteEngine {
	engine := SqlLiteEngine{}
	engine.initEngine(database)
	return &engine
}

func (o *SqlLiteEngine) initEngine(database *config.Database) {
	color.Greenf("初始化数据库。%v", database)
	color.Greenln()
	o.name = DefaultSqlEngineName
	o.url = database.Url
}

func (o *SqlLiteEngine) Name() string {
	return o.name
}

func (o *SqlLiteEngine) Start() {
	db, err := gorm.Open(sqlite.Open(o.url), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	o.db = db
	o.autoMigrate() //创建更新表
}

func (o *SqlLiteEngine) Stop() {
}

func (o *SqlLiteEngine) autoMigrate() {
	err := o.db.AutoMigrate(&AutoUser{})
	if err != nil {
		panic(err)
	}
}

// Raw 执行SQL
func (o *SqlLiteEngine) Raw(sql string, values any, dest any) {
	o.db.Raw(sql, values).Scan(&dest)
}

// Insert 插入对象
func (o *SqlLiteEngine) Insert(values any) {
	o.db.Create(values)
}

type AutoUser struct {
	ID   uint   // Standard field for the primary key
	Name string // A regular string field
}
