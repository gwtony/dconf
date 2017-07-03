package member

import (
	"encoding/json"
	dconf "github.com/gwtony/dconf/dcserver/handler"
	"github.com/gwtony/dconf/test_tools/utils"
)

func MemberAdd1() string {
	mm := dconf.MemberMessage{}
	mm.Service = "testservice"
	mm.Ip = "1.1.1.1"
	mmv, _ := json.Marshal(mm)
	status, _, err := utils.SendRequest(dconf.MEMBER_ADD_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func MemberAdd2() string {
	mm := dconf.MemberMessage{}
	mm.Service = "testservice"
	mm.Ip = "1.1.1.2"
	mmv, _ := json.Marshal(mm)
	status, _, err := utils.SendRequest(dconf.MEMBER_ADD_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func MemberRead() string {
	mm := dconf.MemberMessage{}

	mm.Service = "testservice"
	mm.Group = "default"
	mmv, _ := json.Marshal(mm)

	status, resp, err := utils.SendRequest(dconf.MEMBER_READ_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	mr := dconf.MemberReply{}
	err = json.Unmarshal(resp, &mr)
	if err != nil {
		return "unmarshal failed"
	}
	if mr.Result[0].Group != "default" {
		return "read group is invalid"
	}
	if mr.Result[0].Ip[0] != "1.1.1.1" {
		return "read ip is invalid"
	}

	return ""
}

func MemberReadTestgroup() string {
	mm := dconf.MemberMessage{}

	mm.Service = "testservice"
	mm.Group = "testgroup"
	mmv, _ := json.Marshal(mm)

	status, resp, err := utils.SendRequest(dconf.MEMBER_READ_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	mr := dconf.MemberReply{}
	err = json.Unmarshal(resp, &mr)
	if err != nil {
		return "unmarshal failed"
	}
	if mr.Result[0].Group != "testgroup" {
		return "read group is invalid"
	}
	if mr.Result[0].Ip[0] != "1.1.1.1" {
		return "read ip is invalid"
	}

	return ""
}
func MemberReadNone() string {
	mm := dconf.MemberMessage{}

	mm.Service = "testservice"
	mm.Group = "default"
	mmv, _ := json.Marshal(mm)

	status, _, err := utils.SendRequest(dconf.MEMBER_READ_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 204 {
		return "status is not 204"
	}

	return ""
}

func MemberDelete1() string {
	mm := dconf.MemberMessage{}
	mm.Service = "testservice"
	mm.Group = "default"
	mm.Ip = "1.1.1.1"
	mmv, _ := json.Marshal(mm)
	status, _, err := utils.SendRequest(dconf.MEMBER_DELETE_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}

func MemberDelete2() string {
	mm := dconf.MemberMessage{}
	mm.Service = "testservice"
	mm.Group = "default"
	mm.Ip = "1.1.1.2"
	mmv, _ := json.Marshal(mm)
	status, _, err := utils.SendRequest(dconf.MEMBER_DELETE_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}

func MemberMove() string {
	mm := dconf.MemberMoveMessage{}
	mm.Service = "testservice"
	mm.From = "default"
	mm.To = "testgroup"
	mm.Ip = "1.1.1.1"
	mmv, _ := json.Marshal(mm)
	status, _, err := utils.SendRequest(dconf.MEMBER_MOVE_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}
func MemberMoveBack() string {
	mm := dconf.MemberMoveMessage{}
	mm.Service = "testservice"
	mm.From = "testgroup"
	mm.To = "default"
	mm.Ip = "1.1.1.1"
	mmv, _ := json.Marshal(mm)
	status, _, err := utils.SendRequest(dconf.MEMBER_MOVE_LOC, string(mmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}
