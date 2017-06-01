package handler
import (
	//"fmt"
	//"time"
	//"strings"
	//"strconv"
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
		h.log.Error("Method invalid, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error("Read from body failed, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := &ServiceMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from body failed, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Service add request: (%s), client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" {
		h.log.Error("Service not exist, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Service not exist"), errors.BadRequestError, h.log)
		return
	}
	if data.Description == "" {
		h.log.Error("Read from body failed, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Description not exist"), errors.BadRequestError, h.log)
		return
	}

	if !IsAdmin(r) {
		h.log.Debug("return a failed")
		api.ReturnError(r, w, errors.Jerror("Authentication failed"), errors.UnauthorizedError, h.log)
		return
	}
	//check service exists
	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	h.log.Debug("key is ", key)
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Service add check serivce key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg != nil {
		h.log.Error("Service add service: %s exists", data.Service)
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
		h.log.Error("Service add set key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Add service to backend failed"), errors.BadGatewayError, h.log)
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
		h.log.Error("Method invalid, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error("Read from body failed, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := &ServiceMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from body failed, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Service delete request: (%s), client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" {
		h.log.Error("Service not exist, client: %s", r.RemoteAddr)
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
		h.log.Error("Service delete check group failed: key is %s", key)
		api.ReturnError(r, w, errors.Jerror("Delete service check group failed"), errors.BadGatewayError, h.log)
		return
	}
	if msg != nil && len(msg) > 2 { //group is more than 'default' and 'all'
		h.log.Info("Service delete exists more group")
		api.ReturnError(r, w, errors.Jerror("Delete service need delete group first"), errors.NotAcceptableError, h.log)
		return
	}

	//unset token
	key = h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Service delete unset service meta key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service meta to backend failed"), errors.BadGatewayError, h.log)
		return
	}
	h.log.Debug("Unset service meta done")

	//unset srv view
	key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service
	err = h.eh.UnSetDir(key)
	if err != nil {
		h.log.Error("Service delete unset service view key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service view to backend failed"), errors.BadGatewayError, h.log)
		return
	}
	h.log.Debug("Unset service view done")

	//unset group view
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service
	err = h.eh.UnSetDir(key)
	if err != nil {
		h.log.Error("Service delete unset group view key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service group view to backend failed"), errors.BadGatewayError, h.log)
		return
	}
	h.log.Debug("Unset group view done")

	//unset group meta
	key = h.eh.root + ETCD_GROUP_META + "/" + data.Service
	err = h.eh.UnSetDir(key)
	if err != nil {
		h.log.Error("Service delete unset group meta key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service group meta to backend failed"), errors.BadGatewayError, h.log)
		return
	}
	h.log.Debug("Unset group meta done")

	api.ReturnResponse(r, w, "", h.log)
}

func (h *ServiceReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.log.Error("Method invalid, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error("Read from body failed, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := &ServiceMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from body failed, client: %s", r.RemoteAddr)
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Service read request: (%s), client: %s", string(result), r.RemoteAddr)

	//check args
	if data.Service == "" {
		h.log.Error("Service not exist, client: %s", r.RemoteAddr)
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
		h.log.Error("Service read get servicer meta key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Read service to backend failed"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Error("Service read got no msg key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Read service to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, string(msg.Value), h.log)
}
