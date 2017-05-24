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
	h.log.Info("Service add request: (%s), client: %s", data, r.RemoteAddr)

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
		api.ReturnError(r, w, errors.Jerror("Authentication failed"), errors.UnauthorizedError, h.log)
		return
	}
	//check service exists
	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
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

	sm := &ServiceMessage{Service: data.Service, Description: data.Description, Token: data.Token}
	smv, _ := json.Marshal(sm)

	//set service meta
	key = h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	err = h.eh.Set(key, string(smv))
	if err != nil {
		h.log.Error("Service add set key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Add service to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	////set srv view
	//key = h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	////key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service
	//_, err := h.eh.Set(key, token) //TODO: set dir and default group
	//if err != nil {
	//	h.log.Error("Service add set key %s failed", key)
	//	api.ReturnError(r, w, errors.Jerror("Init service view to backend failed"), errors.BadGatewayError, h.log)
	//	return
	//}
	sr := &ServiceReply{}
	sr.Token = token
	srv, _ := json.Marshal(sr)

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
	h.log.Info("Service add request: (%s), client: %s", data, r.RemoteAddr)

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
		api.ReturnError(r, w, errors.Jerror("Authentication failed"), errors.UnauthorizedError, h.log)
		return
	}

	//unset token
	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Service delete unset key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	//TODO: unset or ...
	//unset srv view
	key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service
	err = h.eh.UnSetDir(key) //TODO: unset dir recurisve
	if err != nil {
		h.log.Error("Service delete unset key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Delete service view to backend failed"), errors.BadGatewayError, h.log)
		return
	}

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
	h.log.Info("Service read request: (%s), client: %s", data, r.RemoteAddr)

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
		h.log.Error("Service read get key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Read service to backend failed"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Error("Service read get key %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Read service to backend failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, string(msg.Value), h.log)
}
