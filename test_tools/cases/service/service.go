package service

import (
	"encoding/json"
	dconf "github.com/gwtony/dconf/dcserver/handler"
	"github.com/gwtony/dconf/test_tools/utils"
)

func ServiceAdd() string {
	sm := dconf.ServiceMessage{}
	sm.Service = "testservice"
	sm.Description = "a test service description"
	smv, _ := json.Marshal(sm)
	status, resp, err := utils.SendRequest(dconf.SERVICE_ADD_LOC, string(smv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	sr := dconf.ServiceReply{}
	_ = json.Unmarshal(resp, &sr)
	if len(sr.Token) > 0 {
		return ""
	}

	return "token invalid"
}

func ServiceRead() string {
	sm := dconf.ServiceMessage{}
	sm.Service = "testservice"
	smv, _ := json.Marshal(sm)
	status, resp, err := utils.SendRequest(dconf.SERVICE_READ_LOC, string(smv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	if resp == nil {
		return "reponse is nil"
	}

	sr := dconf.ServiceMessage{}
	err = json.Unmarshal(resp, &sr)
	if err != nil {
		return "Unmarshal failed"
	}
	if sr.Service == "" {
		return "Service invalid"
	}
	if sr.Description == "" {
		return "Description invalid"
	}

	//fmt.Println(sr)
	return ""
}

func ServiceDelete() string {
	sm := dconf.ServiceMessage{}
	sm.Service = "testservice"
	smv, _ := json.Marshal(sm)
	status, _, err := utils.SendRequest(dconf.SERVICE_DELETE_LOC, string(smv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func ServiceClean() string {
	sm := dconf.ServiceMessage{}
	sm.Service = "testservice"
	smv, _ := json.Marshal(sm)
	_, _, err := utils.SendRequest(dconf.SERVICE_DELETE_LOC, string(smv), "", true)
	if err != nil {
		return "request failed"
	}
	return ""
}
