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

type Item  map[string]string
type Map   map[string]Item
type GItem []string
type GMap  map[string]GItem
type IItem map[string]bool
type IMap  map[string]IItem

func DoRender(eh *EtcdHandler, rm *RenderMessage, log log.Log) (rmsg, errmsg string, err error) {
	var ip string
	log.Info("In do render")

	now := int(time.Now().UnixNano())
	kversion := rm.Key
	if rm.Key == "*" {
		kversion = "wildcard"
	}
	version := fmt.Sprintf("%d_%s_%s_%s", now, rm.Service, kversion, rm.Tag)
	log.Debug("Version is %s", version)

	//check wildcard
	kwildcard := false
	gwildcard := false
	if rm.Key == "*" {
		kwildcard = true
	}
	if rm.Group == "*" {
		gwildcard = true
	}

	key := ""
	group := ""

	if gwildcard { // group is "*"
		key = eh.root + ETCD_SERVICE_VIEW + "/" + rm.Service + "/"
	} else if kwildcard { // group is not "*", key is "*"
		key = eh.root + ETCD_SERVICE_VIEW + "/" + rm.Service + "/" + rm.Group + "/"
	} else { // group is not "*", key is not "*"
		key = eh.root + ETCD_SERVICE_VIEW + "/" + rm.Service + "/" + rm.Group + "/" + rm.Key
	}

	smsg, err := eh.GetWithPrefix(key)
	if err != nil {
		log.Error("Render do get key: %s failed", key)
		//api.ReturnError(r, w, errors.Jerror("Cannot check host with backend"), errors.BadGatewayError, h.log)
		return "", "Cannot check host with backend", errors.BadGatewayError
	}
	if len(smsg) == 0 {
		log.Debug("Render do not exist any key, delete in host view")
	}

	// fill group key map
	gkmap := Map{}
	gki := Item{}
	for _, m := range smsg {
		arr := strings.Split(string(m.Key), "/")
		// key cannot contain "/"
		key = arr[len(arr) - 1]

		if group != arr[4] {
			group = arr[4]
			gki = Item{}
			//log.Debug("add %s: %s to gkmap", group, gki)
			gkmap[group] = gki
		}

		gki[key] = m.Value
	}
	//log.Debug("gkmap is ", gkmap)

	if gwildcard {
		key = eh.root + ETCD_GROUP_VIEW + "/" + rm.Service + "/"
	} else {
		key = eh.root + ETCD_GROUP_VIEW + "/" + rm.Service + "/" + rm.Group
	}
	gmsg, err := eh.GetWithPrefix(key)
	if err != nil {
		log.Error("Render do get with prefix key: %s failed", key)
		return "", "Render do get ip failed", errors.BadGatewayError
		//api.ReturnError(r, w, errors.Jerror("Render do get ip failed"), errors.BadGatewayError, h.log)
		//return
	}
	if gmsg == nil || len(gmsg) == 0 {
		log.Info("Render do no ip in group: %s", rm.Group)
		return "", "Render do no ip in group", errors.NoContentError
		//api.ReturnError(r, w, errors.Jerror("Render do no ip in group"), errors.NoContentError, h.log)
		//return
	}

	group = ""
	gimap := GMap{}
	gii := GItem{}
	//log.Debug("gip len is %d", len(gmsg))
	iplist := make([]string, 0, len(gmsg) / 2) //len is ip num * 2
	// fill group ip map
	for _, m := range gmsg {
		arr := strings.Split(string(m.Key), "/")
		ip = arr[len(arr) - 1]
		//ignore group 'all'
		if gwildcard {
			if arr[4] == "all" {
				iplist = append(iplist, ip)
				continue
			}
		} else {
			iplist = append(iplist, ip)
		}

		if group != arr[4] {
			gii = GItem{}
		}
		group = arr[4]
		gii = append(gii, ip)
		gimap[group] = gii
	}
	//log.Debug("gimap is ", gimap)
	//log.Debug("iplist is ", iplist)

	//fill ip key map
	ikmap := IMap{}
	for _, ip := range iplist {
		if kwildcard {
		//TEST: if gwildcard || kwildcard {
			key = eh.root + ETCD_HOST_VIEW + "/" + ip + "/" + rm.Service + "/"
		} else {
			key = eh.root + ETCD_HOST_VIEW + "/" + ip + "/" + rm.Service + "/" + rm.Key
		}
		//log.Debug("XXX get prefix key: %s", key)
		hmsg, err := eh.GetWithPrefix(key)
		if err != nil {
			log.Error("Render do get keys from host view failed")
			return "", "Render do get keys failed", errors.BadGatewayError
			//api.ReturnError(r, w, errors.Jerror("Render do get keys failed"), errors.BadGatewayError, h.log)
			//return
		}
		//log.Debug("XXX get prefix key result len is: %d", len(hmsg))

		iki := IItem{}
		for _, m := range hmsg {
			arr := strings.Split(string(m.Key), "/")
			// key cannot contain "/"
			key = arr[len(arr) - 1]

			iki[key] = true
			ikmap[ip] = iki
			//log.Debug("XXX add iki: %s to ikmap[%s]", iki, ip)
		}
	}
	//log.Debug("ikmap is ", ikmap)
	log.Info("do render to unset")

	// Check to unset
	for g, ipl := range gimap {
		for _, ip := range ipl {
			ipv, ok := ikmap[ip]
			if !ok {
				//log.Debug("Render do continue, group: %s", g)
				continue
			}
			for ipk, _ := range ipv {
				if _, ok := gkmap[g][ipk]; ok {
					continue
				}

				//log.Debug("Render do wildcard unset key: %s, host: %s in host view", ipk, ip)
				key = eh.root + ETCD_HOST_VIEW + "/" + ip + "/" + rm.Service + "/" + ipk
				err = eh.UnSet(key)
				if err != nil {
					log.Error("Render do wilcard unset key: %s failed", key)
					return "", "Render do wildcard unset failed", errors.BadGatewayError
					//api.ReturnError(r, w, errors.Jerror("Render do wildcard unset failed"), errors.BadGatewayError, h.log)
					//return
				}
			}
		}
	}

	log.Info("do render to set")
	// Check to set
	for g, _ := range gkmap {
		ipl, ok := gimap[g]
		if !ok {
			//log.Debug("Render do continue, group: %s", g)
			continue
		}
		for _, ip := range ipl {
			for k, v := range gkmap[g] {
				if !kwildcard && k != rm.Key {
					//log.Debug("key: %s not match continue", k)
					continue
				}
				//log.Debug("Render do wildcard set key: %s, host: %s in host view", k, ip)
				key = eh.root + ETCD_HOST_VIEW + "/" + ip + "/" + rm.Service + "/" + k
				err = eh.Set(key, v)
				if err != nil {
					log.Error("Render do wilcard set key: %s failed", key)
					return "", "Render do wildcard set failed", errors.BadGatewayError
					//api.ReturnError(r, w, errors.Jerror("Render do wildcard set failed"), errors.BadGatewayError, h.log)
					//return
				}
			}
		}
	}

	rr := &RenderReply{Version: version}
	rrv, _ := json.Marshal(rr)
	return string(rrv), "", nil
}

func (h *RenderDoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//var ip string

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

	rmsg, emsg, err := DoRender(h.eh, data, h.log)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror(emsg), err, h.log)
		return
	}

	api.ReturnResponse(r, w, rmsg, h.log)
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
			h.log.Error("Render read get with prefix key: %s failed", key)
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
			h.log.Error("Render read get key: %s failed", key)
			api.ReturnError(r, w, errors.Jerror("Cannot check ip with backend"), errors.BadGatewayError, h.log)
			return
		}
		if msg == nil {
			h.log.Info("Render read key: %s not exist", data.Key)
			api.ReturnError(r, w, errors.Jerror("Render read key not exist"), errors.NoContentError, h.log)
			return
		}
		arr := strings.Split(msg.Key, "/")
		// key cannot contain "/"
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
	msg, err := h.eh.Get(key)
	if err != nil {
		h.log.Error("Render delete get key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check key with backend"), errors.BadGatewayError, h.log)
		return
	}
	if msg == nil {
		h.log.Error("Render delete not found key: %s", key)
		api.ReturnError(r, w, errors.Jerror("Not found key"), errors.NoContentError, h.log)
		return
	}

	err = h.eh.UnSet(key)
	if err != nil {
		h.log.Error("Render delete unset key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot delete key with backend"), errors.BadGatewayError, h.log)
		return
	}

	api.ReturnResponse(r, w, "", h.log)
}

