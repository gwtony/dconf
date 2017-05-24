package handler
import (
	"fmt"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/gwtony/gapi/log"
	//"github.com/gwtony/gapi/utils"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
)

type RenderDoHandler struct {
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
	h.log.Info("Render do request: (%s) from client: %s", data, r.RemoteAddr)

	if data.Service == "" {
		api.ReturnError(r, w, errors.Jerror("Service invalid"), errors.BadRequestError, h.log)
		return
	}
	if data.Group == "" {
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
			h.log.Error("Render do token not match")
			api.ReturnError(r, w, errors.Jerror("Authentication failed"), err, h.log)
			return
		}
	}

	now := int(time.Now().Unix())
	version := fmt.Sprintf("%d_%s_%s", now, data.Service, data.Tag)
	h.log.Debug("Version is %s", version)

	//check service
	key := h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/" + data.Key
	kmsg, err := h.eh.Get(key)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Cannot check host with backend"), errors.BadGatewayError, h.log)
		return
	}
	if kmsg == nil {
		api.ReturnError(r, w, errors.Jerror("Render do key not exist"), errors.NoContentError, h.log)
		return
	}

	key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/" + data.Group
	gmsg, err := h.eh.GetWithPrefix(key)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Render do get ip faild"), errors.BadGatewayError, h.log)
		return
	}
	if gmsg == nil {
		api.ReturnError(r, w, errors.Jerror("Render do no ip in group"), errors.NoContentError, h.log)
		return
	}

	for _, m := range gmsg {
		arr := strings.Split(string(m.Key), "/")
		ip = arr[len(arr) - 1]
		key = h.eh.root + ETCD_HOST_VIEW + "/" + ip + "/ " + data.Service + "/" + data.Key
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

