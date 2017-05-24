package main

import (
	"fmt"
	"time"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/dconf/dcagent/handler"
)

func main() {
	err := api.Init("dcagent.conf")
	if err != nil {
		fmt.Println("[Error] Init api failed")
		return
	}
	config := api.GetConfig()
	log := api.GetLog()

	err = handler.InitContext(config, log)
	if err != nil {
		fmt.Println("[Error] Init dcagent failed")
		time.Sleep(time.Second)
		return
	}

	api.Run()
}
