package handler

type ConfigMessage struct {
	Service string `json:"service"`
	Group string `json:"group"`
	Key string `json:"key"`
	Value string `json:"value"`
}

type ConfigReplyMessage struct {
	Result string `json:"result"`
}

type WatchMessage struct {
	Type string `json:"type"`
	Key string
	Value string
}
