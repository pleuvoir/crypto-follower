package engine

import (
	"crypto-follower/core/event"
	"testing"
)

func TestNewMainEngine(t *testing.T) {

	mainEngine := NewMainEngine(event.NewEventEngine())

	mainEngine.Start()
	var result AutoUser
	mainEngine.Raw("select * from auto_users where id=?", 1, &result)

	t.Logf("%v", result)
}
