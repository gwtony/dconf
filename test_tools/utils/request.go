package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"fmt"
	dconf "github.com/gwtony/dconf/dcserver/handler"
)

func SendRequest(loc, data, token string, admin bool) (int, []byte, error) {
	pdata := bytes.NewBufferString(data)
	req, err := http.NewRequest("POST", DconfUrl + loc, pdata)
	if admin {
		req.Header.Add(dconf.ADMIN_TOKEN_HEADER, AdminToken)
	} else {
		req.Header.Add(dconf.USER_TOKEN_HEADER, token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error is ", err)
		return 0, nil, err
	}
	defer resp.Body.Close()
	//fmt.Println("response code is", resp.StatusCode)

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		fmt.Println("read from body failed")
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, body, nil
}
