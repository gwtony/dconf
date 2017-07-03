package handler
import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/utils"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
)

type ServiceAddHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type ServiceDeleteHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type ServiceReadHandler struct {
	eh  *EtcdHandler
	log log.Log
}

func (h *ServiceAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &ServiceMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Service add request: (%s), client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Description == "" {
		api.ReturnError(r, w, errors.Jerror("Description invalid"), errors.BadRequestError, h.log)
		return
	}
	if utils.StringToBytes("_")[0] == data.Service[0] {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}

	if !IsAdmin(r) {
		api.ReturnError(r, w, errors.Jerror("Authentication failed"), errors.UnauthorizedError, h.log)
		return
	}

	//check service exists
	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Service add get key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg != nil {
		h.log.Info("Service add service: %s already exists", data.Service)
		api.ReturnError(r, w, errors.Jerror("Service exists"), errors.ConflictError, h.log)
		return
	}

	//generate token
	token, _ := utils.NewToken()

	sm := &ServiceMessage{Service: data.Service, Description: data.Description, Token: token}
	smv, _ := json.Marshal(sm)

	//set service meta
	key = h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	err = h.eh.Set(key, string(smv))
	if err != nil {
		h.log.Error("Service add set service meta key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Add service to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	//set group default meta
	gm := &GroupMessage{
		Service: data.Service,
		Group: "default",
		Description: "default group",
	}
	gmv, _ := json.Marshal(gm)
	key = h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/default"
	err = h.eh.Set(key, string(gmv))
	if err != nil {
		h.log.Error("Service add set default group failed")
		api.ReturnError(r, w, errors.Jerror("Add default group to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	sr := &ServiceReply{}
	sr.Token = token
	srv, _ := json.Marshal(sr)
	h.log.Debug("token is %s", token)

	api.ReturnResponse(r, w, string(srv), h.log)
}

func (h *ServiceDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &ServiceMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Service delete request: (%s), client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("Service not exist"), errors.BadRequestError, h.log)
		return
	}

	if !IsAdmin(r) {
		api.ReturnError(r, w, errors.Jerror("Authentication failed"), errors.UnauthorizedError, h.log)
		return
	}

	// Cannot delete if group is more than 'default' and 'all'
	key := h.eh.root + ETCD_GROUP_META + "/" + data.Service
	msg, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Service delete get with prefix key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service check group failed"), errors.BadGatewayError, h.log)
		return
	}
	if msg != nil && len(msg) > 2 { //group is more than 'default' and 'all'
		h.log.Info("Service delete exists more group")
		api.ReturnError(r, w, errors.Jerror("Need to delete group first"), errors.NotAcceptableError, h.log)
		return
	}

	//unset token
	key = h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Service delete unset service meta key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service meta to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	//unset srv view
	key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service
	err = h.eh.UnSetDir(key)
	if err != nil {
		h.log.Error("Service delete unset service view key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service view to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	//unset group view
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service
	err = h.eh.UnSetDir(key)
	if err != nil {
		h.log.Error("Service delete unset group view key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service group view to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	//unset group meta
	key = h.eh.root + ETCD_GROUP_META + "/" + data.Service
	err = h.eh.UnSetDir(key)
	if err != nil {
		h.log.Error("Service delete unset group meta key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service group meta to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *ServiceReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &ServiceMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Service read request: (%s), client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("Service not exist"), errors.BadRequestError, h.log)
		return
	}

	if !IsAdmin(r) {
		api.ReturnError(r, w, errors.Jerror("Authentication failed"), errors.UnauthorizedError, h.log)
		return
	}

	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Service read get servicer meta key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Read service to backend failed"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Service read got no msg in key: %s", key)
		api.ReturnError(r, w, errors.Jerror("Read service to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, string(msg.Value), h.log)
}
