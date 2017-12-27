package handler

import (
	"time"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
)

type TagAddHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type TagDeleteHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type TagInfoHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type TagReadHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type TagApplyHandler struct {
	eh  *EtcdHandler
	log log.Log
}

func (h *TagAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &TagRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Tag add request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Tag == "" || strings.Contains(data.Tag, "/") {
		api.ReturnError(r, w, errors.Jerror("Tag invalid"), errors.BadRequestError, h.log)
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

	// check group exist
	key := h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.Group
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Tag add get group meta key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Tag add group: %s not exists", data.Group)
		api.ReturnError(r, w, errors.Jerror("Group to not exist"), errors.NoContentError, h.log)
		return
	}

	// check tag conflict
	key = h.eh.root + ETCD_TAG_VIEW + "/" + data.Service + "/" + data.Tag
	msg, err = h.eh.Get(key)
	if err != nil {
		h.log.Error("Tag add get tag key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check tag with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg != nil {
		h.log.Info("Tag add check tag: %s already exists", data.Tag)
		api.ReturnError(r, w, errors.Jerror("Tag exist"), errors.ConflictError, h.log)
		return
	}

	// check tag num limit
	key = h.eh.root + ETCD_TAG_VIEW + "/" + data.Service + "/"
	msgs, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Tag add get tag key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check tag with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msgs != nil && len(msgs) >= TAG_LIMIT_MAX {
		h.log.Info("Tag add check tag reach limit")
		//TODO: return some status code
		api.ReturnError(r, w, errors.Jerror("Tag size limited"), errors.BadRequestError, h.log)
		return
	}

	// get keys from service view
	key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/"
	msgs, err = h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Tag add get with prefix key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msgs == nil {
		h.log.Info("Tag add key: %s not exists", key)
		api.ReturnError(r, w, errors.Jerror("Group not exist"), errors.NoContentError, h.log)
		return
	}

	// store tag 
	tm := TagMeta{Service: data.Service, Group: data.Group, Tag: data.Tag}
	tm.Kv = make(map[string]string, len(msgs))

	for _, m := range msgs {
		arr := strings.Split(string(m.Key), "/")
		// key cannot contain "/"
		mkey := arr[len(arr) - 1]
		tm.Kv[mkey] = m.Value
	}

	tmv, _ := json.Marshal(tm)

	key = h.eh.root + ETCD_TAG_VIEW + "/" + data.Service + "/" + data.Tag
	err = h.eh.Set(key, string(tmv))
	if err != nil {
		h.log.Error("Tag add set tag key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Set tag to backend failed"), err, h.log)
		return
	}
	//TODO: store tag to git

	api.ReturnResponse(r, w, "", h.log)
}

func (h *TagDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &TagRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Tag delete request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
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

	key := h.eh.root + ETCD_TAG_VIEW + "/" + data.Service + "/" + data.Tag
	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Tag delete set tag key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Unset tag to backend failed"), err, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *TagReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &TagRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Tag read request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
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

	key := h.eh.root + ETCD_TAG_VIEW + "/" + data.Service + "/" + data.Tag
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Tag read get tag key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check tag with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Tag read check tag: %s not exists", data.Tag)
		api.ReturnError(r, w, errors.Jerror("Tag not exist"), errors.NoContentError, h.log)
		return
	}

	api.ReturnResponse(r, w, string(msg.Value), h.log)
}

func (h *TagInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &TagInfoRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Tag add request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
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

	key := h.eh.root + ETCD_TAG_VIEW + "/" + data.Service + "/" + data.Group + "/"
	msg, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Tag info get tag key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check tag with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil || len(msg) == 0 {
		h.log.Info("Tag read check tag: (service:%s, group:%s) not exists", data.Service, data.Group)
		api.ReturnError(r, w, errors.Jerror("Tag not exist"), errors.NoContentError, h.log)
		return
	}
	tir := TagInfoReply{}
	tir.Tags = make([]string, 0, len(msg))
	for _, m := range msg {
		arr := strings.Split(string(m.Key), "/")
		mkey := arr[len(arr) - 1]
		tir.Tags = append(tir.Tags, mkey)
	}

	tirv, _ := json.Marshal(tir)

	api.ReturnResponse(r, w, string(tirv), h.log)
}

func (h *TagApplyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &TagRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Tag apply request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Tag == "" || strings.Contains(data.Tag, "/") {
		api.ReturnError(r, w, errors.Jerror("Tag invalid"), errors.BadRequestError, h.log)
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

	key := h.eh.root + ETCD_TAG_VIEW + "/" + data.Service + "/" + data.Tag
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Tag apply get tag key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check tag with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Tag apply check tag: %s not exists", data.Tag)
		api.ReturnError(r, w, errors.Jerror("Tag not exist"), errors.NoContentError, h.log)
		return
	}

	tm := TagMeta{}
	err = json.Unmarshal([]byte(msg.Value), &tm)
	if err != nil {
		h.log.Info("Tag apply check tag: %s not exists", data.Tag)
		api.ReturnError(r, w, errors.Jerror("Tag not exist"), errors.InternalServerError, h.log)
		return
	}
	service := tm.Service
	group := tm.Group
	key = h.eh.root + ETCD_SERVICE_VIEW + "/" + service + "/" + group + "/"
	msgs, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Tag apply get tag key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check tag with backend"), errors.BadGatewayError, h.log)
		return
	}

	// fill old key map
	var okeymap map[string]bool
	if msgs == nil || len(msgs) <= 0 {
		okeymap = make(map[string]bool, 1)
	} else {
		okeymap = make(map[string]bool, len(msgs))
		for _, m := range msgs {
			arr := strings.Split(string(m.Key), "/")
			// key cannot contain "/"
			mkey := arr[len(arr) - 1]
			okeymap[mkey] = true
		}
	}

	// delete old key
	for k, _ := range okeymap {
		if _, ok := tm.Kv[k]; !ok {
			cm := &ConfigMessage{Key: k, Service: service, Group: group}
			h.log.Debug("Tag apply delete config: %s/%s/%s", service, group, k)
			DeleteConfig(h.eh, cm, h.log)
		}
	}

	// update or add key
	for tmk, tmv := range tm.Kv {
		if _, ok := okeymap[tmk]; ok {
			cm := &ConfigMessage{Key: tmk, Value: tmv, Service: service, Group: group}
			h.log.Debug("Tag apply update config: %s/%s/%s", service, group, tmk)
			UpdateConfig(h.eh, cm, h.log)
		} else {
			cm := &ConfigMessage{Key: tmk, Value: tmv, Service: service, Group: group}
			h.log.Debug("Tag apply add config: %s/%s/%s", service, group, tmk)
			AddConfig(h.eh, cm, h.log)
		}
	}

	ts := strconv.Itoa(int(time.Now().Unix()))
	rm := &RenderMessage{Service: service, Group: group, Key: "*", Tag: "tag_" + ts}
	rmsg, emsg, err := DoRender(h.eh, rm, h.log)
	h.log.Info("Tag apply render ret: ", rmsg, emsg, err)
	//TODO: return error

	api.ReturnResponse(r, w, "", h.log)
}
