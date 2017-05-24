package handler

const (
	// VERSION version
	VERSION                 = "0.1 alpha"

	API_CONTENT_HEADER      = "application/json;charset=utf-8"

	ETCD_EVENT_PUT          = "PUT"
	ETCD_EVENT_DELETE       = "DELETE"

	DEFAULT_ADMIN_TOKEN     = "LCONF_TOKEN"
	ADMIN_TOKEN_HEADER      = "Admin-Token"
	USER_TOKEN_HEADER       = "Conf-Token"

	DCONF_LOC               = "/dconf"
	DEFAULT_DICT_SIZE       = 10000

	ETCD_USER               = "dconf_user"
	ETCD_PASSWORD           = "dconf_password"

	SET                     = iota
	UNSET
	GET
	GET_PREFIX

	CONTENT_HEADER          = "Content-Type"
	ETCD_CONTENT_HEADER     = "application/x-www-form-urlencoded"
	DEFAULT_ETCD_PORT       = "2379"
	DEFAULT_ETCD_TIMEOUT    = 3
	DEFAULT_ETCD_ROOT       = "/dconf"

	//ETCD_ARGS
	ETCD_V2_PREFIX          = "/v2/keys"
	ETCD_RECURSIEVE_ARGS    = "?recursive=true"
	ETCD_DIR_RECU_ARGS      = "?dir=true&recursive=true"
	ETCD_DIR_PADDING        = "true"

	// ETCD table schema
	ETCD_SERVICE_VIEW       = "/svc_view"
	ETCD_GROUP_VIEW         = "/group_view"
	ETCD_HOST_VIEW          = "/host_view"
	ETCD_TAG_VIEW           = "/tag_view"

	ETCD_SERVICE_META       = "/srv_meta"
	ETCD_GROUP_META         = "/group_meta"

	ETCD_GROUP_DEFAULT      = "/default"
	ETCD_GROUP_ALL          = "/all"

	ETCD_IP_PADDING         = "alive"


	//http location for config
	CONFIG_GET_LOC          = "/config/get"
)
