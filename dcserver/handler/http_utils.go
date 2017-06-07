package handler
import (
	"encoding/json"
	"net/http"
	"github.com/gwtony/gapi/errors"
)

func IsAdmin(r *http.Request) bool {
	if r.Header[ADMIN_TOKEN_HEADER] != nil {
		if r.Header[ADMIN_TOKEN_HEADER][0] == AdminToken {
			return true
		}
	}
	return false
}

func CheckToken(r *http.Request, eh *EtcdHandler, service string) (bool, error){
	if r.Header[USER_TOKEN_HEADER] == nil {
		eh.log.Debug("No user token header")
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
		eh.log.Info("Token key: %s not exists", key)
		return false, errors.NotAcceptableError
	}
	sm := &ServiceMessage{}
	err = json.Unmarshal([]byte(msg.Value), &sm)
	if err != nil {
		eh.log.Error("Parse server meta failed: %s", key)
		return false, errors.InternalServerError
	}

	if sm.Token != token {
		eh.log.Debug("Token not match")
		return false, errors.UnauthorizedError
	}

	return true, nil
}
