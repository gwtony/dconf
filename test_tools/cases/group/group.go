package group

import (
	//"fmt"
	"encoding/json"
	dconf "github.com/gwtony/dconf/dcserver/handler"
	"github.com/gwtony/dconf/test_tools/utils"
)

func GroupAdd() string {
	mm := dconf.GroupMessage{}
	mm.Service = "testservice"
	mm.Group = "testgroup"
	mm.Description = "thisisgroupdescription"
	mmv, _ := json.Marshal(mm)
	status, _, err := utils.SendRequest(dconf.GROUP_ADD_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func GroupRead() string {
	mm := dconf.GroupMessage{}

	mm.Service = "testservice"
	mm.Group = "testgroup"
	mmv, _ := json.Marshal(mm)

	status, resp, err := utils.SendRequest(dconf.GROUP_READ_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	rmm := dconf.GroupMessage{}
	err = json.Unmarshal(resp, &rmm)
	if err != nil {
		return "unmarshal failed"
	}
	if rmm.Group != "testgroup" {
		return "read group is invalid"
	}

	return ""
}

func GroupReadUpdate() string {
	mm := dconf.GroupMessage{}

	mm.Service = "testservice"
	mm.Group = "testgroup"
	mmv, _ := json.Marshal(mm)

	status, resp, err := utils.SendRequest(dconf.GROUP_READ_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	rmm := dconf.GroupMessage{}
	err = json.Unmarshal(resp, &rmm)
	if err != nil {
		return "unmarshal failed"
	}
	if rmm.Group != "testgroup" {
		return "read group is invalid"
	}
	if rmm.Description != "update" {
		return "check update  is invalid"
	}

	return ""
}

func GroupReadNone() string {
	mm := dconf.GroupMessage{}

	mm.Service = "testservice"
	mm.Group = "testgroup"
	mmv, _ := json.Marshal(mm)

	status, _, err := utils.SendRequest(dconf.GROUP_READ_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 204 {
		return "status is not 204"
	}

	return ""
}

func GroupDelete() string {
	mm := dconf.GroupMessage{}
	mm.Service = "testservice"
	mm.Group = "testgroup"
	mmv, _ := json.Marshal(mm)
	status, _, err := utils.SendRequest(dconf.GROUP_DELETE_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}

func GroupList() string {
	mm := dconf.GroupMessage{}
	mm.Service = "testservice"
	mmv, _ := json.Marshal(mm)
	status, resp, err := utils.SendRequest(dconf.GROUP_LIST_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	rmm := dconf.GroupReply{}
	err = json.Unmarshal(resp, &rmm)
	if err != nil {
		return "unmarshal failed"
	}
	if len(rmm.Result) < 0 {
		return "read no result"
	}

	if rmm.Result[0].Group != "testgroup" {
		return "group wrong"
	}


	return ""
}

func GroupUpdate() string {
	mm := dconf.GroupMessage{}
	mm.Service = "testservice"
	mm.Group = "testgroup"
	mm.Description = "update"
	mmv, _ := json.Marshal(mm)
	status, _, err := utils.SendRequest(dconf.GROUP_UPDATE_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

