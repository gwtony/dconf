package handler
import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/gwtony/gapi/log"
	//"github.com/gwtony/gapi/utils"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
)

type ConfigGetHandler struct {
	eh  *EtcdHandler
	log log.Log
}

func (h *ConfigGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	h.log.Info("Config get request: (%s) from client: %s", data, r.RemoteAddr)

	//check args
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

	key := data.Service + "/" + data.Group + "/" + data.Key

	res := GetConfig(key)

	crm := &ConfigReplyMessage{}
	crm.Result = res

	crmv, _ := json.Marshal(crm)

	api.ReturnResponse(r, w, string(crmv), h.log)
}
