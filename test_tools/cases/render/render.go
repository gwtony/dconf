package render

import (
	"encoding/json"
	dconf "github.com/gwtony/dconf/dcserver/handler"
	"github.com/gwtony/dconf/test_tools/utils"
)

func RenderDo() string {
	cm := dconf.RenderMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config"
	cm.Tag = "test_tag"
	cmv, _ := json.Marshal(cm)
	status, resp, err := utils.SendRequest(dconf.RENDER_DO_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	rr := &dconf.RenderReply{}
	if err := json.Unmarshal(resp, &rr); err != nil {
		return "unmarsh failed"
	}

	if rr.Version == "" {
		return "bad version"
	}

	return ""
}

func RenderDo2() string {
	cm := dconf.RenderMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config2"
	cm.Tag = "test_tag"
	cmv, _ := json.Marshal(cm)
	status, resp, err := utils.SendRequest(dconf.RENDER_DO_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	rr := &dconf.RenderReply{}
	if err := json.Unmarshal(resp, &rr); err != nil {
		return "unmarsh failed"
	}

	if rr.Version == "" {
		return "bad version"
	}

	return ""
}

func RenderRead() string {
	cm := dconf.RenderReadMessage{}
	cm.Service = "testservice"
	cm.Ip = "1.1.1.1"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)
	status, resp, err := utils.SendRequest(dconf.RENDER_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	if resp == nil {
		return "return body failed"
	}
	rr := dconf.RenderReadReply{}
	err = json.Unmarshal(resp, &rr)
	if err != nil {
		return "unmarshal failed"
	}
	if len(rr.Result) == 0 {
		return "len invalid"
	}
	if rr.Result[0].Key != "test_config" {
		return "key not match"
	}
	if rr.Result[0].Value != "test_config_value" {
		return "value not match"
	}
	return ""
}

func RenderReadWildcard() string {
	cm := dconf.RenderReadMessage{}
	cm.Service = "testservice"
	cm.Ip = "1.1.1.1"
	cm.Key = "*"
	cmv, _ := json.Marshal(cm)
	status, resp, err := utils.SendRequest(dconf.RENDER_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	if resp == nil {
		return "return body failed"
	}
	rr := dconf.RenderReadReply{}
	err = json.Unmarshal(resp, &rr)
	if err != nil {
		return "unmarshal failed"
	}
	if len(rr.Result) != 2 {
		return "len invalid"
	}
	if rr.Result[0].Key != "test_config" {
		return "key not match"
	}
	if rr.Result[0].Value != "test_config_value" {
		return "value not match"
	}
	if rr.Result[1].Key != "test_config2" {
		return "key not match"
	}
	if rr.Result[1].Value != "test_config_value2" {
		return "value not match"
	}
	return ""
}

func RenderDelete() string {
	cm := dconf.RenderDeleteMessage{}
	cm.Service = "testservice"
	cm.Ip = "1.1.1.1"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)
	status, _, err := utils.SendRequest(dconf.RENDER_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func RenderDeleteWildcard() string {
	cm := dconf.RenderDeleteMessage{}
	cm.Service = "testservice"
	cm.Ip = "1.1.1.1"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)
	status, _, err := utils.SendRequest(dconf.RENDER_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	cm.Key = "test_config2"
	cmv, _ = json.Marshal(cm)
	status, _, err = utils.SendRequest(dconf.RENDER_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func RenderReadNone() string {
	cm := dconf.RenderReadMessage{}
	cm.Service = "testservice"
	cm.Ip = "1.1.1.1"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)
	status, _, err := utils.SendRequest(dconf.RENDER_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 204 {
		return "status is not 200"
	}
	return ""
}
