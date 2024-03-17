package engine

import (
	"crypto-follower/core/config"
	"crypto-follower/core/helper"
	"testing"
)

func TestNewDbEngine(t *testing.T) {
	engine := NewDbEngine(&config.Database{Name: DefaultSqlEngineName, Url: "test.db"})

	engine.Start()

	user := AutoUser{Name: "Jinzhu"}

	engine.Insert(&user) // 通过数据的指针来创建

	t.Log(user.ID)

	var result2 AutoUser
	engine.Raw("select * from auto_users where id=?", 1, &result2)

	t.Logf("%v", result2)

	path, _ := helper.CurrentPath()
	t.Log(path)
}
