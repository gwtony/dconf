package handler

import (
	//"sync"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/config"
)

var AdminToken string

// InitContext inits dconf context
func InitContext(conf *config.Config, log log.Log) error {
	cf := &LConfConfig{}
	err := cf.ParseConfig(conf)
	if err != nil {
		log.Error("Lconf parse config failed")
		return err
	}
	AdminToken = cf.adminToken

	//mc, err := InitMysqlContext(cf.maddr, cf.dbname, cf.dbuser, cf.dbpwd, log)
	//if err != nil {
	//	log.Error("Dcron init mysql context failed")
	//	return err
	//}


	eh := InitEtcdHandler(cf.eaddr, cf.eto, cf.euser, cf.epwd, cf.eauthEnable, cf.eRoot, log)

	apiLoc := cf.apiLoc

	api.AddHttpHandler(apiLoc + SERVICE_ADD_LOC, &ServiceAddHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + SERVICE_DELETE_LOC, &ServiceDeleteHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + SERVICE_READ_LOC, &ServiceReadHandler{eh: eh, log: log})

	api.AddHttpHandler(apiLoc + GROUP_ADD_LOC, &GroupAddHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + GROUP_DELETE_LOC, &GroupDeleteHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + GROUP_UPDATE_LOC, &GroupUpdateHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + GROUP_READ_LOC, &GroupReadHandler{eh: eh, log: log})

	api.AddHttpHandler(apiLoc + MEMBER_ADD_LOC, &MemberAddHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + MEMBER_DELETE_LOC, &MemberDeleteHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + MEMBER_MOVE_LOC, &MemberMoveHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + MEMBER_READ_LOC, &MemberReadHandler{eh: eh, log: log})

	api.AddHttpHandler(apiLoc + CONFIG_ADD_LOC, &ConfigAddHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + CONFIG_DELETE_LOC, &ConfigDeleteHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + CONFIG_READ_LOC, &ConfigReadHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + CONFIG_UPDATE_LOC, &ConfigUpdateHandler{eh: eh, log: log})
	api.AddHttpHandler(apiLoc + CONFIG_COPY_LOC, &ConfigCopyHandler{eh: eh, log: log})

	api.AddHttpHandler(apiLoc + RENDER_DO_LOC, &RenderDoHandler{eh: eh, log: log})

	//ch.Run()

	return nil
}
