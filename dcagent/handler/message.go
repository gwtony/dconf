package handler

type ConfigMessage struct {
	Service string `json:"service"`
	Ip      string `json:"ip"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

type ConfigMeta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ConfigReplyMessage struct {
	Result []*ConfigMeta `json:"result"`
}

type WatchMessage struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
