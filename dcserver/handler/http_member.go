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

type MemberAddHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type MemberDeleteHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type MemberReadHandler struct {
	eh  *EtcdHandler
	log log.Log
}

type MemberMoveHandler struct {
	eh  *EtcdHandler
	log log.Log
}

// add to all and default group
func (h *MemberAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &MemberMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Member add request: (%s) from client: %s", string(result), r.RemoteAddr)

	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("No Service"), errors.BadRequestError, h.log)
		return
	}
	if data.Ip == "" {
		api.ReturnError(r, w, errors.Jerror("No ip"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)
	if !admin {
		api.ReturnError(r, w, errors.Jerror("Authentication failed"), errors.UnauthorizedError, h.log)
		return
	}

	//check service
	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	msg, err := h.eh.Get(key)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		api.ReturnError(r, w, errors.Jerror("Service not exist"), errors.NoContentError, h.log)
		return
	}

	//check ip exists in all
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + ETCD_GROUP_ALL + "/" + data.Ip
	msg, err = h.eh.Get(key)
	if err != nil {
		h.log.Error("Member add check ip: %s failed", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Cannot check host with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg != nil {
		h.log.Info("Member add ip: %s already exists", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Member add ip exist"), errors.ConflictError, h.log)
		return
	}

	//set to all group
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + ETCD_GROUP_ALL + "/" + data.Ip
	value := ETCD_IP_PADDING
	err = h.eh.Set(key, value)
	if err != nil {
		h.log.Error("Member add set ip: %s to all failed", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Member add ip failed"), errors.BadGatewayError, h.log)
		return
	}

	//set to default group
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + ETCD_GROUP_DEFAULT + "/" + data.Ip
	value = ETCD_IP_PADDING
	err = h.eh.Set(key, value)
	if err != nil {
		h.log.Error("Member add set ip: %s to default failed", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Member add ip failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *MemberDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &MemberMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Member delete request: (%s) from client: %s", string(result), r.RemoteAddr)

	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("No service"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" {
		api.ReturnError(r, w, errors.Jerror("No group"), errors.BadRequestError, h.log)
		return
	}
	if data.Ip == "" {
		api.ReturnError(r, w, errors.Jerror("No ip"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "all" {
		api.ReturnError(r, w, errors.Jerror("Group name invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)
	if !admin {
		h.log.Error("Member delete check admin failed")
		api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
		return
	}

	//check ip exists in all
	key := h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + ETCD_GROUP_ALL + "/" + data.Ip
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Member delete get ip: %s from group all failed", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Cannot check host with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Member delete get group all found no ip")
		api.ReturnError(r, w, errors.Jerror("Member delete ip not exist"), errors.NoContentError, h.log)
		return
	}

	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/" + data.Group + "/" + data.Ip
	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Member delete unset ip: %s from group failed", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Member delete ip from group failed"), errors.BadGatewayError, h.log)
		return
	}

	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + ETCD_GROUP_ALL + "/" + data.Ip
	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Member delete unset ip: %s from all failed", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Member delete ip from all failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

func (h *MemberReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &MemberMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Member read request: (%s) from client: %s", string(result), r.RemoteAddr)

	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("No Service"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" {
		api.ReturnError(r, w, errors.Jerror("No group"), errors.BadRequestError, h.log)
		return
	}

	wildcard := false
	wc := "*"
	if data.Group == wc {
		wildcard = true
	}

	// check admin
	admin := IsAdmin(r)
	if !admin {
		ok, err := CheckToken(r, h.eh, data.Service)
		if !ok {
			h.log.Error("Member read token not match")
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	//check service
	key := h.eh.root + ETCD_SERVICE_META + "/" + data.Service
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Member read get service: %s failed", data.Service)
		api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		api.ReturnError(r, w, errors.Jerror("Service not exist"), errors.NoContentError, h.log)
		return
	}

	//check ip exists in all
	if !wildcard {
		key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/" + data.Group
	} else {
		key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/"
	}
	msgarr, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Member read get key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check ip with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msgarr == nil || len(msgarr) == 0 {
		h.log.Info("Member read ip: %s not exist", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Member read ip not exist"), errors.NoContentError, h.log)
		return
	}

	group := ""
	mr := &MemberReply{}
	mm := &MemberMeta{}
	for _, m := range msgarr {
		//this is sorted
		arr := strings.Split(m.Key, "/")
		if arr[len(arr) - 2] == "all" { //skip group "all"
			continue
		}
		if group == "" {
			group = arr[len(arr) - 2]
			mm.Group = group
		}
		if group != arr[len(arr) - 2] {
			mr.Result = append(mr.Result, mm)
			mm = &MemberMeta{}
			mm.Group = arr[len(arr) - 2]
		}

		mm.Ip = append(mm.Ip, arr[len(arr) - 1])

		h.log.Debug("Member read m.key: %s, m.value: %s", m.Key, m.Value)
	}
	mr.Result = append(mr.Result, mm)
	h.log.Debug(mr)
	mrv, _ := json.Marshal(mr)

	api.ReturnResponse(r, w, string(mrv), h.log)
}

func (h *MemberMoveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := &MemberMoveMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Member add request: (%s) from client: %s", string(result), r.RemoteAddr)

	if data.Service == "" || strings.Contains(data.Service, "/") {
		api.ReturnError(r, w, errors.Jerror("Invalid service"), errors.BadRequestError, h.log)
		return
	}
	if data.Ip == "" {
		api.ReturnError(r, w, errors.Jerror("Invalid ip"), errors.BadRequestError, h.log)
		return
	}
	if data.From == "" || strings.Contains(data.From, "/") {
		api.ReturnError(r, w, errors.Jerror("Invalid from"), errors.BadRequestError, h.log)
		return
	}
	if data.To == "" || strings.Contains(data.To, "/") {
		api.ReturnError(r, w, errors.Jerror("Invalid to"), errors.BadRequestError, h.log)
		return
	}
	if data.From == "all" {
		api.ReturnError(r, w, errors.Jerror("Group from invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.To == "all" {
		api.ReturnError(r, w, errors.Jerror("Group to invalid"), errors.BadRequestError, h.log)
		return
	}

	// check admin
	admin := IsAdmin(r)
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
		h.log.Error("Member move get service: %s failed", data.Service)
		api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Info("Member move service: %s not exist", data.Service)
		api.ReturnError(r, w, errors.Jerror("Service not exist"), errors.NoContentError, h.log)
		return
	}

	//check source group exists
	if data.From != "default" {
		key = h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.From
		msg, err = h.eh.Get(key)
		if err != nil {
			h.log.Error("Member move get source group: %s failed", data.From)
			api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
			return
		}
		if msg == nil {
			h.log.Info("Member move source group: %s not exist", data.From)
			api.ReturnError(r, w, errors.Jerror("Group in from not exist"), errors.NoContentError, h.log)
			return
		}
	}

	//chech dest group exists
	if data.To != "default" {
		key = h.eh.root + ETCD_GROUP_META + "/" + data.Service + "/" + data.To
		msg, err = h.eh.Get(key)
		if err != nil {
			h.log.Error("Member move get dest group: %s failed", data.To)
			api.ReturnError(r, w, errors.Jerror("Cannot check service with backend"), errors.BadGatewayError, h.log)
			return
		}
		if msg == nil {
			h.log.Info("Member move dest group: %s not exist", data.To)
			api.ReturnError(r, w, errors.Jerror("Group in to not exist"), errors.NoContentError, h.log)
			return
		}
	}

	//check ip exists in all
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + ETCD_GROUP_ALL + "/" + data.Ip
	msg, err = h.eh.Get(key)
	if err != nil {
		h.log.Error("Member move get ip: %s from group all failed", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Cannot check host with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
			h.log.Info("Member move ip: %s not exist", data.Ip)
		api.ReturnError(r, w, errors.Jerror("Member move ip not exist"), errors.NoContentError, h.log)
		return
	}

	//delete from source
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/" + data.From + "/" + data.Ip
	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Member move unset key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Member move delete from source group failed"), errors.BadGatewayError, h.log)
		return
	}

	//add to dest
	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/" + data.To + "/" + data.Ip
	err = h.eh.Set(key, ETCD_IP_PADDING)
	if err != nil {
		h.log.Error("Member move set key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Member move delete from source group failed"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}
