package main

import (
	"fmt"
	"time"
	"github.com/gwtony/gapi/api"
	"git.lianjia.com/lianjia-sysop/lconf/handler"
)

func main() {
	err := api.Init("lconf.conf")
	if err != nil {
		fmt.Println("[Error] Init api failed")
		return
	}
	config := api.GetConfig()
	log := api.GetLog()

	err = handler.InitContext(config, log)
	if err != nil {
		fmt.Println("[Error] Init lconf failed")
		time.Sleep(time.Second)
		return
	}

	api.Run()
}
