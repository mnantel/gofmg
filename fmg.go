package fmg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FMG struct {
	Username string
	Password string
	IP       string
	Port     string
	sessKey  string
}

type fmgAPI struct {
	Method  string         `json:"method"`
	Params  []fmgAPIParams `json:"params"`
	Session string         `json:"session"`
	ID      string         `json:"id"`
	Ver     string         `json:"ver"`
}

type fmgAPIParams struct {
	URL  string      `json:"url"`
	Data interface{} `json:"data"`
}

type fmgAPIDataLogin struct {
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

type FMGAPIResp struct {
	ID     string `json:"id"`
	Ver    string `json:"ver"`
	Result []struct {
		URL    string `json:"url"`
		Status struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"status"`
		Data []map[string]interface{} `json:"data"`
	} `json:"result"`
	Session string `json:"session"`
}

// Login performs authentication to target FMG and stores resulting session key.
func (fmg *FMG) Login() error {

	FMGLoginData := fmgAPIDataLogin{
		User:   fmg.Username,
		Passwd: fmg.Password,
	}
	var datapost []fmgAPIDataLogin
	datapost = append(datapost, FMGLoginData)

	res, err := fmg.apicall("exec", "/sys/login/user", datapost)
	if err != nil {
		return err
	}

	fmg.sessKey = res.Session
	return nil
}

// Call is an exported function to place an API call to FMG
func (fmg *FMG) Call(method string, url string, data interface{}) (FMGAPIResp, error) {
	var empty FMGAPIResp
	if fmg.sessKey == "" {
		return empty, errors.New("must call Login() first, session key is null")
	}
	res, err := fmg.apicall(method, url, data)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (fmg *FMG) apicall(method string, url string, data interface{}) (FMGAPIResp, error) {

	var resp FMGAPIResp

	client := &http.Client{}

	reqURI := fmt.Sprintf("https://%s:%s/%s", fmg.IP, fmg.Port, "jsonrpc")

	apicallparams := fmgAPIParams{
		URL:  url,
		Data: data,
	}

	apicall := fmgAPI{
		Method:  method,
		ID:      "1",
		Session: fmg.sessKey,
	}

	apicall.Params = append(apicall.Params, apicallparams)
	postdata, err := json.Marshal(apicall)

	req, err := http.NewRequest("POST", reqURI, bytes.NewBuffer(postdata))
	if err != nil {
		return resp, err
	}
	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(resbody, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil

}
