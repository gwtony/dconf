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

type GroupAddHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type GroupDeleteHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type GroupUpdateHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type GroupReadHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type GroupListHandler struct {
	eh  *EtcdHandler
	log log.Log
}

func (h *GroupAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &GroupMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Group add request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Description == "" {
		api.ReturnError(r, w, errors.Jerror("Description invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "default" || data.Group == "all" {
		api.ReturnError(r, w, errors.Jerror("Group name invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)

	// check token
	if !admin {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	//check service
	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Group add check service: %s failed", data.Service)
		api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Group add check service: %s not exist", data.Service)
		api.ReturnError(r, w, errors.Jerror("Service not exist"), errors.NoContentError, h.log)
		return
	}

	//check group meta
	key = h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.Group
	msg, err = h.eh.Get(key)
	if err != nil {
		h.log.Error("Group add check group: %s meta failed", data.Group)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg != nil {
		h.log.Info("Group add check group: %s already exists", data.Group)
		api.ReturnError(r, w, errors.Jerror("Group exist"), errors.ConflictError, h.log)
		return
	}

	//set group meta
	key = h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.Group
	err = h.eh.Set(key, string(result))
	if err != nil {
		h.log.Error("Group add set group: %s meta failed", data.Group)
		api.ReturnError(r, w, errors.Jerror("Add group to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *GroupDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &GroupMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Group delete request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "default" || data.Group == "all" {
		api.ReturnError(r, w, errors.Jerror("Group name invalid"), errors.BadRequestError, h.log)
		return
	}
	// check admin
	admin := IsAdmin(r)

	// check token
	if !admin {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	//check service exists
	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Group delete check service: %s failed", data.Service)
		api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Group delete check service: %s not exist", data.Service)
		api.ReturnError(r, w, errors.Jerror("Service not exist"), errors.NoContentError, h.log)
		return
	}

	//check group exists
	key = h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.Group
	msg, err = h.eh.Get(key)
	if err != nil {
		h.log.Error("Group delete check group: %s failed", data.Group)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Group delete check group %s not exist", data.Group)
		api.ReturnError(r, w, errors.Jerror("Group not exist"), errors.NoContentError, h.log)
		return
	}

	// check group empty
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/" + data.Group
	msga, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Group delete check group get key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check group ip with backend"), errors.BadGatewayError, h.log)
		return
	}
	if len(msga) > 0 {
		h.log.Info("Group delete check group: %s is not empty", data.Group)
		api.ReturnError(r, w, errors.Jerror("Group is not empty"), errors.NoContentError, h.log)
		return
	}

	// delete group view dir
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/" + data.Group
	err = h.eh.UnSetDir(key) //unset dir
	if err != nil {
		h.log.Error("Group delete unsetdir group: %s view failed", data.Group)
		api.ReturnError(r, w, errors.Jerror("Delete group to backend"), errors.BadGatewayError, h.log)
		return
	}

	// delete service view dir
	key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group
	err = h.eh.UnSetDir(key)
	if err != nil {
		h.log.Error("Group delete unsetdir service: %s view failed", data.Service)
		api.ReturnError(r, w, errors.Jerror("Delete group to backend"), errors.BadGatewayError, h.log)
		return
	}

	//unset group meta
	key = h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.Group
	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Group delete unset group: %s meta failed", data.Service)
		api.ReturnError(r, w, errors.Jerror("Add group to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *GroupUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &GroupMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Group update request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" || strings.Contains(data.Group, "/") {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Description == "" {
		api.ReturnError(r, w, errors.Jerror("Description invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "default" || data.Group == "all" {
		api.ReturnError(r, w, errors.Jerror("Group name invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)

	// check token
	if !admin {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	//check service
	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Group update check service: %s meta failed", data.Service)
		api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Group update check service: %s not exist", data.Service)
		api.ReturnError(r, w, errors.Jerror("Service not exist"), errors.NoContentError, h.log)
		return
	}

	//check group
	key = h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.Group
	msg, err = h.eh.Get(key)
	if err != nil {
		h.log.Error("Group update check group: %s meta failed", data.Group)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Group update check group: %s not exist", data.Group)
		api.ReturnError(r, w, errors.Jerror("Group not exist"), errors.NoContentError, h.log)
		return
	}

	//only set group meta
	key = h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.Group
	err = h.eh.Set(key, string(result))
	if err != nil {
		h.log.Error("Group update set group: %s meta failed", data.Group)
		api.ReturnError(r, w, errors.Jerror("Add group to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *GroupReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &GroupMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Group read request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" {
		api.ReturnError(r, w, errors.Jerror("Group invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)

	// check token
	if !admin {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	key := h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.Group
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Group read get group: %s meta failed", data.Group)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Group read check group: %s not exist", data.Group)
		api.ReturnError(r, w, errors.Jerror("Group not exist"), errors.NoContentError, h.log)
		return
	}

	gm := &GroupMessage{}
	err = json.Unmarshal([]byte(msg.Value), &gm)
	if err != nil {
		h.log.Error("Group read unmarshal for: %s/%s failed", data.Service, data.Group)
		api.ReturnError(r, w, errors.Jerror("Unmarshal failed"), errors.InternalServerError, h.log)
		return
	}

	gmv, _ := json.Marshal(gm)

	api.ReturnResponse(r, w, string(gmv), h.log)
}

func (h *GroupListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &GroupMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Group list request: (%s) from client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)

	// check token
	if !admin {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	key := h.eh.root + ETCD_GROUP_META + "/" + data.Service
	msg, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Group list get with prefix service: %s failed", data.Service)
		api.ReturnError(r, w, errors.Jerror("Cannot check group with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Group list no group in service: %s", data.Service)
		api.ReturnError(r, w, errors.Jerror("No group in service"), errors.NoContentError, h.log)
		return
	}

	gr := &GroupReply{}
	gr.Result = make([]*GroupMeta, 0, len(msg))
	for _, m := range msg {
		gm := &GroupMessage{}
		err = json.Unmarshal([]byte(m.Value), &gm)
		if err != nil {
			h.log.Error("Group list unmarshal for: %s/%s failed", data.Service, data.Group)
			api.ReturnError(r, w, errors.Jerror("Unmarshal failed"), errors.InternalServerError, h.log)
			return
		}
		if gm.Group == "all" || gm.Group == "default" { //ignore "all" and "default"
			continue
		}
		gmm := &GroupMeta{Group: gm.Group, Description: gm.Description}
		gr.Result = append(gr.Result, gmm)
	}

	grv, _ := json.Marshal(gr)

	api.ReturnResponse(r, w, string(grv), h.log)
}
