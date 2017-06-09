package handler

import (
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/config"
)

// InitContext inits dconf context
func InitContext(conf *config.Config, log log.Log) error {
	cf := &DCAgentConfig{}
	err := cf.ParseConfig(conf)
	if err != nil {
		log.Error("dcagent parse config failed")
		return err
	}

	//Need not to auth
	eh := InitEtcdHandler(cf.eaddr, cf.eto, "", "", false, cf.eRoot, log)

	dm, err := InitDictManager(cf.localhost, eh, cf.store, log)
	if err != nil {
		log.Error("Init dict manager failed")
		return err
	}
	err = dm.PullAll()
	if err != nil {
		log.Error("Pull config from etcd failed:", err)
	}

	go dm.Run()

	apiLoc := cf.apiLoc
	api.AddHttpHandler(apiLoc + CONFIG_GET_LOC, &ConfigGetHandler{eh: eh, log: log})

	return nil

}
