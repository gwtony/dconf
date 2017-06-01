package handler

import (
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/gwtony/gapi/config"
	"github.com/gwtony/gapi/errors"
)

type DCServerConfig struct {
	adminToken    string        /* admin token */

	eaddr         []string      /* etcd addr */
	eto           time.Duration /* etcd timeout */
	euser         string
	epwd          string
	eauthEnable   bool
	eRoot         string

	apiLoc        string        /* dcserver api location */


	timeout       time.Duration
}

// ParseConfig parses config
func (conf *DCServerConfig) ParseConfig(cf *config.Config) error {
	var err error
	if cf.C == nil {
		return errors.BadConfigError
	}
	eaddrStr, err := cf.C.GetString("dcserver", "etcd_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [dcserver] Read conf: No etcd_addr")
		return err
	}

	if eaddrStr == "" {
		fmt.Fprintln(os.Stderr, "[Error] [dcserver] Read conf: Empty etcd server address")
		return errors.BadConfigError
	}
	eaddr := strings.Split(eaddrStr, ",")
	for i := 0; i < len(eaddr); i++ {
		if eaddr[i] != "" {
			if !strings.Contains(eaddr[i], ":") {
				conf.eaddr = append(conf.eaddr, eaddr[i] + ":" + DEFAULT_ETCD_PORT)
			} else {
				conf.eaddr = append(conf.eaddr, eaddr[i])
			}
		}
	}
	eto, err :=  cf.C.GetInt64("dcserver", "etcd_timeout")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcserver] Read conf: no etcd_timeout, use default timeout:", DEFAULT_ETCD_TIMEOUT)
		eto = DEFAULT_ETCD_TIMEOUT
	}
	conf.eto = time.Duration(eto) * time.Second

	conf.apiLoc, err = cf.C.GetString("dcserver", "api_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcserver] Read conf: No api_location, use default location:", DCSERVER_LOC)
		conf.apiLoc = DCSERVER_LOC
	}

	conf.adminToken, err = cf.C.GetString("dcserver", "admin_token")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcserver] Read conf: No admin_token, use default admin token:", DEFAULT_ADMIN_TOKEN)
		conf.adminToken = DEFAULT_ADMIN_TOKEN
	}

	conf.eRoot, err = cf.C.GetString("dcserver", "etcd_root")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcserver] Read conf: No etcd_root, use default etcd root:", DEFAULT_ETCD_ROOT)
		conf.eRoot = DEFAULT_ETCD_ROOT
	}

	conf.euser, err = cf.C.GetString("dcserver", "etcd_user")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcserver] Read conf: No api_location")
		conf.euser = ""
	}
	conf.epwd, err = cf.C.GetString("dcserver", "etcd_pwd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcserver] Read conf: No etcd_pwd")
		conf.epwd = ""
	}
	conf.eauthEnable = true
	if conf.euser == "" || conf.epwd == "" {
		conf.eauthEnable = false
	}

	return nil
}
