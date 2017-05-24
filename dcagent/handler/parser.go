package handler

func DecodeWatchMessage(otype string, key, value []byte) (*WatchMessage, error) {
	wm := &WatchMessage{}

	wm.Type = otype
	wm.Key = string(key)
	wm.Value = string(value)

	return wm, nil
}
