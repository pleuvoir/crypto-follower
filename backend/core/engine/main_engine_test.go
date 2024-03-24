package engine

import (
	"crypto-follower/core/config"
	"os"
	"testing"
)

func TestNewMainEngine(t *testing.T) {

	_ = os.Setenv("conf.path", "/Users/pleuvoir/dev/space/git/crypto-follower/backend/app.yml")

	conf := config.NewApplicationConf()
	if err := conf.Load(); err != nil {
		panic(err)
	}

	mainEngine := NewMainEngine(NewEventEngine())
	mainEngine.InitEngines(&conf)
	mainEngine.Start()
	var result AutoUser
	mainEngine.Raw("select * from auto_users where id=?", 1, &result)

	t.Logf("%v", result)

	//mainEngine.Stop()
}
