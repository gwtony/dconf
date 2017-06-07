package handler
import (
	"fmt"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
)

type RenderDoHandler struct {
	eh  *EtcdHandler
	log log.Log
}
type RenderReadHandler struct {
	eh  *EtcdHandler
	log log.Log
}
type RenderDeleteHandler struct {
	eh  *EtcdHandler
	log log.Log
}


func (h *RenderDoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ip string

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

	data := &RenderMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Render do request: (%s) from client: %s", string(result), r.RemoteAddr)

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
	if data.Tag == "" {
		api.ReturnError(r, w, errors.Jerror("Tag invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)
	if !admin {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			h.log.Info("Render do token not match")
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	now := int(time.Now().UnixNano())
	version := fmt.Sprintf("%d_%s_%s_%s", now, data.Service, data.Key, data.Tag)
	h.log.Debug("Version is %s", version)

	//check service
	key := h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/" + data.Key
	kmsg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Render do get key: %s faild", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check host with backend"), errors.BadGatewayError, h.log)
		return
	}
	if kmsg == nil {
		h.log.Info("Render do key: %s not exist", key)
		api.ReturnError(r, w, errors.Jerror("Render do key not exist"), errors.NoContentError, h.log)
		return
	}

	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/" + data.Group
	gmsg, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Render do get with prefix key: %s faild", key)
		api.ReturnError(r, w, errors.Jerror("Render do get ip faild"), errors.BadGatewayError, h.log)
		return
	}
	if gmsg == nil {
		h.log.Info("Render do no ip in group: %s", data.Group)
		api.ReturnError(r, w, errors.Jerror("Render do no ip in group"), errors.NoContentError, h.log)
		return
	}

	for _, m := range gmsg {
		arr := strings.Split(string(m.Key), "/")
		// key can contain "/"
		//ip = arr[len(arr) - 1]
		ip = strings.Join(arr[5:], "/")
		key = h.eh.root + ETCD_HOST_VIEW + "/" + ip + "/" + data.Service + "/" + data.Key
		err = h.eh.Set(key, string(kmsg.Value))
		if err != nil {
			h.log.Error("Render do set key: %s faild", key)
			api.ReturnError(r, w, errors.Jerror("Render do process failed"), errors.BadGatewayError, h.log)
			return
		}
	}

	rr := &RenderReply{Version: version}
	rrv, _ := json.Marshal(rr)

	api.ReturnResponse(r, w, string(rrv), h.log)
}

func (h *RenderReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &RenderReadMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Render read request: (%s) from client: %s", string(result), r.RemoteAddr)

	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Ip == "" {
		api.ReturnError(r, w, errors.Jerror("Ip invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Key == "" {
		api.ReturnError(r, w, errors.Jerror("Key invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)
	if !admin {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			h.log.Error("Render read token not match")
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	wildcard := false
	if data.Key == "*" {
		wildcard = true
	}

	key := ""
	rr := &RenderReadReply{}

	if wildcard {
		key = h.eh.root + ETCD_HOST_VIEW + "/" + data.Ip + "/" + data.Service + "/"
		msg, err := h.eh.GetWithPrefix(key)
		if err != nil {
			h.log.Error("Render read get with prefix key: %s faild", key)
			api.ReturnError(r, w, errors.Jerror("Cannot read render with backend"), errors.BadGatewayError, h.log)
			return
		}
		if msg == nil {
			h.log.Info("Render read no service: %s in ip: %s", data.Service, data.Ip)
			api.ReturnError(r, w, errors.Jerror("No service in this ip"), errors.NoContentError, h.log)
			return
		}
		rr.Result = make([]*RenderReadMeta, 0, len(msg))
		for _, m := range msg {
			arr := strings.Split(m.Key, "/")
			// key can contain "/"
			//rrm := &RenderReadMeta{Key: arr[len(arr) - 1], Value: m.Value}
			rrm := &RenderReadMeta{Key: strings.Join(arr[5:], "/"), Value: m.Value}
			rr.Result = append(rr.Result, rrm)
		}

	} else {
		key = h.eh.root + ETCD_HOST_VIEW + "/" + data.Ip + "/" + data.Service + "/" + data.Key
		msg, err := h.eh.Get(key)
		if err != nil {
			h.log.Error("Render read get key: %s faild", key)
			api.ReturnError(r, w, errors.Jerror("Cannot check ip with backend"), errors.BadGatewayError, h.log)
			return
		}
		if msg == nil {
			h.log.Info("Render read key: %s not exist", data.Key)
			api.ReturnError(r, w, errors.Jerror("Render read key not exist"), errors.NoContentError, h.log)
			return
		}
		arr := strings.Split(msg.Key, "/")
		rrm := &RenderReadMeta{Key: arr[len(arr) - 1], Value: msg.Value}
		rr.Result = append(rr.Result, rrm)
	}

	rrv, _ := json.Marshal(rr)

	api.ReturnResponse(r, w, string(rrv), h.log)
}

func (h *RenderDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &RenderReadMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Render delete request: (%s) from client: %s", string(result), r.RemoteAddr)

	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Ip == "" {
		api.ReturnError(r, w, errors.Jerror("Ip invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Key == "" {
		api.ReturnError(r, w, errors.Jerror("Key invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)
	if !admin {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			h.log.Info("Render delete token not match")
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	key := h.eh.root + ETCD_HOST_VIEW + "/" + data.Ip + "/" + data.Service + "/" + data.Key
	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Render delete unset key: %s faild", key)
		api.ReturnError(r, w, errors.Jerror("Cannot delete key with backend"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

