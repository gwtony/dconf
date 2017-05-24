package handler

import (
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/gwtony/gapi/config"
	"github.com/gwtony/gapi/errors"
)

type DCAgentConfig struct {
	adminToken    string        /* admin token */

	eaddr         []string      /* etcd addr */
	eto           time.Duration /* etcd timeout */
	euser         string
	epwd          string
	eauthEnable   bool
	eRoot         string

	apiLoc        string        /* dcagent api location */

	localhost     string

	timeout       time.Duration
}

// ParseConfig parses config
func (conf *DCAgentConfig) ParseConfig(cf *config.Config) error {
	var err error
	if cf.C == nil {
		return errors.BadConfigError
	}
	eaddrStr, err := cf.C.GetString("dcagent", "etcd_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [dcagent] Read conf: No etcd_addr")
		return err
	}

	if eaddrStr == "" {
		fmt.Fprintln(os.Stderr, "[Error] [dcagent] Read conf: Empty etcd server address")
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
	eto, err :=  cf.C.GetInt64("dcagent", "etcd_timeout")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcagent] Read conf: no etcd_timeout, use default timeout:", DEFAULT_ETCD_TIMEOUT)
		eto = DEFAULT_ETCD_TIMEOUT
	}
	conf.eto = time.Duration(eto) * time.Second

	conf.apiLoc, err = cf.C.GetString("dcagent", "api_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcagent] Read conf: No api_location, use default location:", DCONF_LOC)
		conf.apiLoc = DCONF_LOC
	}

	conf.adminToken, err = cf.C.GetString("dcagent", "admin_token")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcagent] Read conf: No admin_token, use default admin token:", DEFAULT_ADMIN_TOKEN)
		conf.adminToken = DEFAULT_ADMIN_TOKEN
	}

	conf.eRoot, err = cf.C.GetString("dcagent", "etcd_root")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcagent] Read conf: No etcd_root, use default etcd root:", DEFAULT_ETCD_ROOT)
		conf.eRoot = DEFAULT_ETCD_ROOT
	}

	conf.euser, err = cf.C.GetString("dcagent", "etcd_user")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcagent] Read conf: No api_location")
		conf.euser = ""
	}
	conf.epwd, err = cf.C.GetString("dcagent", "etcd_pwd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [dcagent] Read conf: No etcd_pwd")
		conf.epwd = ""
	}
	conf.eauthEnable = true
	if conf.euser == "" || conf.epwd == "" {
		conf.eauthEnable = false
	}

	conf.localhost, err = cf.C.GetString("dcagent", "host")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [dcagent] Read conf: No host")
		return err
	}

	return nil
}
