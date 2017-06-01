package handler
import (
	//"time"
	"strings"
	//"io/ioutil"
	"encoding/json"
	"net/http"
	//"github.com/gwtony/gapi/log"
	//"github.com/gwtony/gapi/utils"
	//"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
)

func IsAdmin(r *http.Request) bool {
	if r.Header[ADMIN_TOKEN_HEADER] != nil {
		if (strings.Compare(r.Header[ADMIN_TOKEN_HEADER][0], AdminToken) == 0) {
			return true
		}
	}
	return false
}

func CheckToken(r *http.Request, eh *EtcdHandler, service string) (bool, error){
	if r.Header[USER_TOKEN_HEADER] == nil {
		return false, errors.UnauthorizedError
	}

	token := r.Header[USER_TOKEN_HEADER][0]

	key := eh.root + ETCD_SERVICE_META + "/" + service
	msg, err := eh.Get(key)
	if err != nil {
		eh.log.Error("Get token key %s failed", key)
		return false, errors.BadGatewayError
	}
	if msg == nil {
		eh.log.Error("Token key: %s not exists", key)
		return false, errors.NotAcceptableError
	}
	sm := &ServiceMessage{}
	err = json.Unmarshal([]byte(msg.Value), &sm)
	if err != nil {
		eh.log.Error("Parse server meta failed: %s", key)
		return false, errors.InternalServerError
	}

	if strings.Compare(sm.Token, token) != 0 {
		eh.log.Error("Token not match")
		return false, errors.UnauthorizedError
	}

	return true, nil
}
