package main

import (
	"fmt"
	"time"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/dconf/dcserver/handler"
)

func main() {
	err := api.Init("dcserver.conf")
	if err != nil {
		fmt.Println("[Error] Init api failed")
		return
	}
	config := api.GetConfig()
	log := api.GetLog()

	err = handler.InitContext(config, log)
	if err != nil {
		fmt.Println("[Error] Init dcserver failed")
		time.Sleep(time.Second)
		return
	}

	api.Run()
}
