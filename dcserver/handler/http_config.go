package handler
import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
)

type ConfigAddHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type ConfigDeleteHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type ConfigReadHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type ConfigUpdateHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type ConfigCopyHandler struct {
	eh  *EtcdHandler
	log log.Log
}

func (h *ConfigAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := &ConfigMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Config add request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Key == "" {
		api.ReturnError(r, w, errors.Jerror("Key invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Value == "" {
		api.ReturnError(r, w, errors.Jerror("Value invalid"), errors.BadRequestError, h.log)
		return
	}

	// group "all" should not set config
	if data.Group == "all"{
		api.ReturnError(r, w, errors.Jerror("Group name invalid"), errors.BadRequestError, h.log)
		return
	}
	// config value should less than MAX_VALUE_SIZE(4k)
	if len(data.Value) > MAX_VALUE_SIZE {
		api.ReturnError(r, w, errors.Jerror("Value size too big"), errors.BadRequestError, h.log)
		return
	}

	// check admin and token
	if !IsAdmin(r) {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	// check group
	if data.Group != "default" {
		key := h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.Group
		msg, err := h.eh.Get(key)
		if err != nil {
			h.log.Error("Config add check group get key: %s failed", key)
			api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
			return
		}
		if msg == nil {
			h.log.Info("Config add group: %s not exists", data.Group)
			api.ReturnError(r, w, errors.Jerror("Group not exist"), errors.NoContentError, h.log)
			return
		}
	}

	//set kv
	key := h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/" + data.Key
	err = h.eh.Set(key, string(data.Value))
	if err != nil {
		h.log.Error("Config add set key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Set config to backend failed"), err, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *ConfigDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := &ConfigMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Config delete request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Key == "" {
		api.ReturnError(r, w, errors.Jerror("Key invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin and token
	if !IsAdmin(r) {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	// Need not to check group
	key := h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/" + data.Key
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Config delete get key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Config delete key: %s not exists", key)
		api.ReturnError(r, w, errors.Jerror("Group not exist"), errors.NoContentError, h.log)
		return
	}

	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Config delete unset key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *ConfigReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := &ConfigMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Config read request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Key == "" {
		api.ReturnError(r, w, errors.Jerror("Key invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin and token
	if !IsAdmin(r) {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	wildcard := false
	wc := "*"
	if data.Key == wc {
		wildcard = true
	}
	// Need not to check group
	key := ""
	if wildcard {
		key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/"
	} else {
		key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/" + data.Key
	}
	msg, err := h.eh.GetWithPrefix(key)

	if err != nil {
		h.log.Error("Config read get key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Read key from server failed"), errors.BadGatewayError, h.log)
		return
	}
	if len(msg) == 0 {
		h.log.Info("Config read key: %s not found", key)
		api.ReturnError(r, w, errors.Jerror("Key not found"), errors.NoContentError, h.log)
		return
	}

	cr := &ConfigReply{}
	for _, m := range msg {
		ckv := &ConfigKV{}
		h.log.Debug("key is %s", m.Key)
		arr := strings.Split(m.Key, "/")
		// key cannot contain "/"
		ckv.Key = arr[len(arr) - 1]
		ckv.Value = string(m.Value)
		cr.Result = append(cr.Result, ckv)
	}
	crv, _ := json.Marshal(cr)

	api.ReturnResponse(r, w, string(crv), h.log)
}

func (h *ConfigUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := &ConfigMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Config update request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Key == "" {
		api.ReturnError(r, w, errors.Jerror("Key invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Value == "" {
		api.ReturnError(r, w, errors.Jerror("Value invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin and token
	if !IsAdmin(r) {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	// check key only
	key := h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/" + data.Key
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Config update get key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Config update key: %s not exists", key)
		api.ReturnError(r, w, errors.Jerror("Group not exist"), errors.NoContentError, h.log)
		return
	}

	//set kv
	key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/" + data.Key
	err = h.eh.Set(key, string(data.Value))
	if err != nil {
		h.log.Error("Config update set key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Set config to backend failed"), err, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *ConfigCopyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := &ConfigCopyMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Config copy request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.From == "" || strings.Contains(data.From, "/") {
		api.ReturnError(r, w, errors.Jerror("From invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.To == "" || strings.Contains(data.To, "/") {
		api.ReturnError(r, w, errors.Jerror("To invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Key == "" {
		api.ReturnError(r, w, errors.Jerror("Key invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin and token
	if !IsAdmin(r) {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	//check dest group
	key := h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.To
	gmsg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Config copy get group meta key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if gmsg == nil {
		h.log.Info("Config copy dest group: %s not exists", data.To)
		api.ReturnError(r, w, errors.Jerror("Group to not exist"), errors.NoContentError, h.log)
		return
	}

	wildcard := false
	wc := "*"
	if data.Key == wc {
		wildcard = true
	}

	if wildcard {
		key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.From + "/"
	} else {
		key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.From + "/" + data.Key
	}

	// get key from 
	msg, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Config copy get with prefix key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Config copy key: %s not exists", key)
		api.ReturnError(r, w, errors.Jerror("Group not exist"), errors.NoContentError, h.log)
		return
	}

	for _, m := range msg {
		arr := strings.Split(string(m.Key), "/")
		mkey := arr[len(arr) - 1]
		key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.To + "/" + mkey
		err = h.eh.Set(key, string(m.Value))
		if err != nil {
			h.log.Error("Config copy set key: %s failed", key)
			api.ReturnError(r, w, errors.Jerror("Set config to backend failed"), err, h.log)
			return
		}
	}

	api.ReturnResponse(r, w, "", h.log)
}
