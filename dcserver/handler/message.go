package handler

type ServiceMessage struct {
	Service     string `json:"service"`
	Description string `json:"description"`
	Token       string `json:"token"`
}

type GroupMessage struct {
	Service     string `json:"service"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

type MemberMessage struct {
	Service string `json:"service"`
	Group   string `json:"group"`
	Ip      string `json:"ip"`
}

type MemberMoveMessage struct {
	Service string `json:"service"`
	From    string `json:"from"`
	To      string `json:"to"`
	Ip      string `json:"ip"`
}

type ConfigMessage struct {
	Service string `json:"service"`
	Group   string `json:"group"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

type ConfigCopyMessage struct {
	Service string `json:"service"`
	From    string `json:"from"`
	To      string `json:"to"`
	Key     string `json:"key"`
}

type RenderMessage struct {
	Service string `json:"service"`
	Group   string `json:"group"`
	Key     string `json:"key"`
	Tag     string `json:"tag"`
}

type RenderDeleteMessage RenderReadMessage
type RenderReadMessage struct {
	Service string `json:"service"`
	Ip      string `json:"ip"`
	Key     string `json:"key"`
}

type RenderReadMeta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RenderReadReply struct {
	Result []*RenderReadMeta `json:"result"`
}

type ServiceReply struct {
	Token string `json:"token"`
}

type GroupReply struct {
	Result []*GroupMeta `json:"result"`
}

type GroupMeta struct {
	Group       string `json:"group"`
	Description string `json:"description"`
}

type MemberMeta struct {
	Group string   `json:"group"`
	Ip    []string `json:ip`
}
type MemberReply struct {
	Result []*MemberMeta `json:"result"`
}

type ConfigReply struct {
	Result []*ConfigKV `json:"result"`
}

type ConfigKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}


type RenderReply struct {
	Version string `json:"version"`
}

type TagRequest struct {
	Tag     string `json:"tag"`
	Service string `json:"service"`
	Group   string `json:"group"`
}

type TagMeta struct {
	Service string            `json:"service"`
	Group   string            `json:"group"`
	Tag     string            `json:"tag"`
	Kv      map[string]string `json:"kv"`
}

type TagInfoRequest struct {
	Service string   `json:"service"`
	Group   string   `json:"group"`
}

type TagInfoReply struct {
	Tags []string `json:"tags"`
}
