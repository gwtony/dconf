package config

import (
	"encoding/json"
	dconf "github.com/gwtony/dconf/dcserver/handler"
	"github.com/gwtony/dconf/test_tools/utils"
)

func ConfigAdd() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config"
	cm.Value = "test_config_value"
	cmv, _ := json.Marshal(cm)
	status, _, err := utils.SendRequest(dconf.CONFIG_ADD_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func ConfigAdd2() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config2"
	cm.Value = "test_config_value2"
	cmv, _ := json.Marshal(cm)
	status, _, err := utils.SendRequest(dconf.CONFIG_ADD_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func ConfigAddSlash() string {
	cm := lconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config/a/slash/key"
	cm.Value = "test_config_value"
	cmv, _ := json.Marshal(cm)
	status, _, err := utils.SendRequest(lconf.CONFIG_ADD_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func ConfigAddGroup() string {
	cm := lconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "testgroup"
	cm.Key = "test_config_group"
	cm.Value = "test_config_value_group"
	cmv, _ := json.Marshal(cm)
	status, _, err := utils.SendRequest(lconf.CONFIG_ADD_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func ConfigUpdate() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config"
	cm.Value = "update"
	cmv, _ := json.Marshal(cm)
	status, _, err := utils.SendRequest(dconf.CONFIG_ADD_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func ConfigRead() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)

	status, resp, err := utils.SendRequest(dconf.CONFIG_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	cr := dconf.ConfigReply{}
	err = json.Unmarshal(resp, &cr)
	if err != nil {
		return "unmarshal failed"
	}
	if cr.Result[0].Key != "test_config" {
		return "read group is invalid"
	}
	if cr.Result[0].Value != "test_config_value" {
		return "read ip is invalid"
	}

	return ""
}

func ConfigReadSlash() string {
	cm := lconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config/a/slash/key"
	cmv, _ := json.Marshal(cm)

	status, resp, err := utils.SendRequest(lconf.CONFIG_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	cr := lconf.ConfigReply{}
	err = json.Unmarshal(resp, &cr)
	if err != nil {
		return "unmarshal failed"
	}
	if cr.Result[0].Key != "test_config/a/slash/key" {
		return "read group is invalid"
	}
	if cr.Result[0].Value != "test_config_value" {
		return "read ip is invalid"
	}

	return ""
}

func ConfigReadUpdate() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)

	status, resp, err := utils.SendRequest(dconf.CONFIG_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	cr := dconf.ConfigReply{}
	err = json.Unmarshal(resp, &cr)
	if err != nil {
		return "unmarshal failed"
	}
	if cr.Result[0].Key != "test_config" {
		return "read group is invalid"
	}
	if cr.Result[0].Value != "update" {
		return "read ip is invalid"
	}

	return ""
}

func ConfigReadNone() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)

	status, _, err := utils.SendRequest(dconf.CONFIG_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 204 {
		return "status is not 204"
	}

	return ""
}

func ConfigReadNoneSlash() string {
	cm := lconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config/a/slash/key"
	cmv, _ := json.Marshal(cm)

	status, _, err := utils.SendRequest(lconf.CONFIG_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 204 {
		return "status is not 204"
	}

	return ""
}

func ConfigDelete() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)

	status, _, err := utils.SendRequest(dconf.CONFIG_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}

func ConfigDelete2() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config2"
	cmv, _ := json.Marshal(cm)

	status, _, err := utils.SendRequest(dconf.CONFIG_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}

func ConfigDeleteSlash() string {
	cm := lconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "default"
	cm.Key = "test_config/a/slash/key"
	cmv, _ := json.Marshal(cm)

	status, _, err := utils.SendRequest(lconf.CONFIG_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}

func ConfigDeleteGroup() string {
	cm := lconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "testgroup"
	cm.Key = "test_config_group"
	cmv, _ := json.Marshal(cm)

	status, _, err := utils.SendRequest(lconf.CONFIG_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}

func ConfigCopy() string {
	ccm := dconf.ConfigCopyMessage{}
	ccm.Service = "testservice"
	ccm.From = "default"
	ccm.To = "testgroup"
	ccm.Key = "test_config"
	cmv, _ := json.Marshal(ccm)

	status, _, err := utils.SendRequest(dconf.CONFIG_COPY_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}
func ConfigDeleteCopy() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "testgroup"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)

	status, _, err := utils.SendRequest(dconf.CONFIG_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}

func ConfigReadCopy() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "testgroup"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)

	status, resp, err := utils.SendRequest(dconf.CONFIG_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	cr := dconf.ConfigReply{}
	err = json.Unmarshal(resp, &cr)
	if err != nil {
		return "unmarshal failed"
	}
	if cr.Result[0].Key != "test_config" {
		return "read group is invalid"
	}
	if cr.Result[0].Value != "test_config_value" {
		return "read ip is invalid"
	}

	return ""
}
func ConfigCopyWildcard() string {
	ccm := dconf.ConfigCopyMessage{}
	ccm.Service = "testservice"
	ccm.From = "default"
	ccm.To = "testgroup"
	ccm.Key = "*"
	cmv, _ := json.Marshal(ccm)

	status, _, err := utils.SendRequest(dconf.CONFIG_COPY_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	return ""
}

func ConfigDeleteCopyWildcard() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "testgroup"
	cm.Key = "test_config"
	cmv, _ := json.Marshal(cm)

	status, _, err := utils.SendRequest(dconf.CONFIG_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	cm.Key = "test_config2"
	cmv, _ = json.Marshal(cm)

	status, _, err = utils.SendRequest(dconf.CONFIG_DELETE_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}
	return ""
}

func ConfigReadCopyWildcard() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "testgroup"
	cm.Key = "*"
	cmv, _ := json.Marshal(cm)

	status, resp, err := utils.SendRequest(dconf.CONFIG_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 200 {
		return "status is not 200"
	}

	cr := dconf.ConfigReply{}
	err = json.Unmarshal(resp, &cr)
	if err != nil {
		return "unmarshal failed"
	}
	if len(cr.Result) != 2 {
		return "read copy wildcard len wrong"
	}
	if cr.Result[0].Key != "test_config" && cr.Result[1].Key != "test_config" {
		return "read key is invalid"
	}
	if cr.Result[0].Value != "test_config_value" && cr.Result[1].Value != "test_config_value2" {
		return "read value is invalid"
	}

	return ""
}

func ConfigReadNoneCopyWildcard() string {
	cm := dconf.ConfigMessage{}
	cm.Service = "testservice"
	cm.Group = "testgroup"
	cm.Key = "*"
	cmv, _ := json.Marshal(cm)

	status, _, err := utils.SendRequest(dconf.CONFIG_READ_LOC, string(cmv), "", true)
	if err != nil {
		return "request failed"
	}
	if status != 204 {
		return "status is not 204"
	}

	return ""
}

