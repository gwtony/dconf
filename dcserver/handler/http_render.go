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
	kversion := data.Key
	if data.Key == "*" {
		kversion = "wildcard"
	}
	version := fmt.Sprintf("%d_%s_%s_%s", now, data.Service, kversion, data.Tag)
	h.log.Debug("Version is %s", version)

	//check wildcard
	kwildcard := false
	gwildcard := false
	if data.Key == "*" {
		kwildcard = true
	}
	if data.Group == "*" {
		gwildcard = true
	}

	key := ""
	group := ""

	if gwildcard { // group is "*"
		key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/"
	} else if kwildcard { // group is not "*", key is "*"
		key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/"
	} else { // group is not "*", key is not "*"
		key = h.eh.root + ETCD_SERVICE_VIEW + "/" + data.Service + "/" + data.Group + "/" + data.Key
	}

	smsg, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Render do get key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Cannot check host with backend"), errors.BadGatewayError, h.log)
		return
	}
	if len(smsg) == 0 {
		h.log.Debug("Render do not exist any key, delete in host view")
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
			h.log.Debug("add %s: %s to gkmap", group, gki)
			gkmap[group] = gki
		}

		gki[key] = m.Value
	}
	h.log.Debug("gkmap is ", gkmap)

	if gwildcard {
		key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/"
	} else {
		key = h.eh.root + ETCD_GROUP_VIEW + "/" + data.Service + "/" + data.Group
	}
	gmsg, err := h.eh.GetWithPrefix(key)
	if err != nil {
		h.log.Error("Render do get with prefix key: %s failed", key)
		api.ReturnError(r, w, errors.Jerror("Render do get ip failed"), errors.BadGatewayError, h.log)
		return
	}
	if gmsg == nil || len(gmsg) == 0 {
		h.log.Info("Render do no ip in group: %s", data.Group)
		api.ReturnError(r, w, errors.Jerror("Render do no ip in group"), errors.NoContentError, h.log)
		return
	}

	group = ""
	gimap := GMap{}
	gii := GItem{}
	h.log.Debug("gip len is %d", len(gmsg))
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
	h.log.Debug("gimap is ", gimap)
	h.log.Debug("iplist is ", iplist)

	//fill ip key map
	ikmap := IMap{}
	for _, ip := range iplist {
		if kwildcard {
			key = h.eh.root + ETCD_HOST_VIEW + "/" + ip + "/" + data.Service + "/"
		} else {
			key = h.eh.root + ETCD_HOST_VIEW + "/" + ip + "/" + data.Service + "/" + data.Key
		}
		h.log.Debug("Get prefix key: %s", key)
		hmsg, err := h.eh.GetWithPrefix(key)
		if err != nil {
			h.log.Error("Render do get keys from host view failed")
			api.ReturnError(r, w, errors.Jerror("Render do get keys failed"), errors.BadGatewayError, h.log)
			return
		}
		h.log.Debug("Get prefix key result len is: %d", len(hmsg))

		iki := IItem{}
		for _, m := range hmsg {
			arr := strings.Split(string(m.Key), "/")
			// key cannot contain "/"
			key = arr[len(arr) - 1]

			iki[key] = true
			ikmap[ip] = iki
			h.log.Debug("Add iki: %s to ikmap[%s]", iki, ip)
		}
	}
	h.log.Debug("ikmap is ", ikmap)

	// Check to unset
	for g, ipl := range gimap {
		for _, ip := range ipl {
			ipv, ok := ikmap[ip]
			if !ok {
				h.log.Debug("Render do continue, group: %s", g)
				continue
			}
			for ipk, _ := range ipv {
				if _, ok := gkmap[g][ipk]; ok {
					continue
				}

				h.log.Debug("Render do wildcard unset key: %s, host: %s in host view", ipk, ip)
				key = h.eh.root + ETCD_HOST_VIEW + "/" + ip + "/" + data.Service + "/" + ipk
				err = h.eh.UnSet(key)
				if err != nil {
					h.log.Error("Render do wilcard unset key: %s failed", key)
					api.ReturnError(r, w, errors.Jerror("Render do wildcard unset failed"), errors.BadGatewayError, h.log)
					return
				}
			}
		}
	}

	// Check to set
	for g, _ := range gkmap {
		ipl, ok := gimap[g]
		if !ok {
			h.log.Debug("Render do continue, group: %s", g)
			continue
		}
		for _, ip := range ipl {
			for k, v := range gkmap[g] {
				if !kwildcard && k != data.Key {
					h.log.Debug("key: %s not match continue", k)
					continue
				}
				h.log.Debug("Render do wildcard set key: %s, host: %s in host view", k, ip)
				key = h.eh.root + ETCD_HOST_VIEW + "/" + ip + "/" + data.Service + "/" + k
				err = h.eh.Set(key, v)
				if err != nil {
					h.log.Error("Render do wilcard set key: %s failed", key)
					api.ReturnError(r, w, errors.Jerror("Render do wildcard set failed"), errors.BadGatewayError, h.log)
					return
				}
			}
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
			// key cannot contain "/"
			rrm := &RenderReadMeta{Key: arr[len(arr) - 1], Value: m.Value}
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

