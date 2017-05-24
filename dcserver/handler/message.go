package handler

type ServiceMessage struct {
	Service string `json:"service"`
	Description string `json:"description"`
	Token string `json:"token"`
}

type GroupMessage struct {
	Service string `json:"service"`
	Group string `json:"group"`
	Description string `json:"description"`
}

type MemberMessage struct {
	Service string `json:"service"`
	Group string `json:"group"`
	Ip string `json:"ip"`
}

type MemberMoveMessage struct {
	Service string `json:"service"`
	From string `json:"from"`
	To string `json:"to"`
	Ip string `json:"ip"`
}

type ConfigMessage struct {
	Service string `json:"service"`
	Group string `json:"group"`
	Key string `json:"key"`
	Value string `json:"value"`
}

type ConfigCopyMessage struct {
	Service string `json:"service"`
	From string `json:"from"`
	To string `json:"to"`
	Key string `json:"key"`
}

type RenderMessage struct {
	Service string `json:"service"`
	Group string `json:"group"`
	Key string `json:"key"`
	Tag string `json:"tag"`
}

type ServiceReply struct {
	Token string `json:"token"`
}

type GroupReply struct {
	Result []GroupMeta `json:"result"`
}

type GroupMeta struct {
	Group string `json:"group"`
	Description string `json:"description"`
}

type MemberMeta struct {
	Group string `json:"group"`
	Ip []string `json:ip`
}
type MemberReply struct {
	Result []MemberMeta `json:"result"`
}

type ConfigReply struct {
	Result []*ConfigKV `json:"result"`
}

type ConfigKV struct {
	Key string `json:"key"`
	Value string `json:"value"`
}


type RenderReply struct {
	Version string `json:"version"`
}

