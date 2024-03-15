package main

import (
	"crypto-follower/bootstrap"
	"crypto-follower/restful"
)

func main() {
	bootstrap.Init()
	restful.NewServer().ServerStartedListener(nil).Start()
}
