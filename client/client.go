package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"genshin-sign-helper/util"
	"genshin-sign-helper/util/constant"
	log "genshin-sign-helper/util/logger"
)

type GenshinClient struct {
	*http.Client
}

func NewGenshinClient() (g *GenshinClient) {
	g = &GenshinClient{
		Client: &http.Client{},
	}
	return
}

func (g *GenshinClient) NewGenshinRequest(method, cookie, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to new http request:%v", err)
	}

	req.Header.Add("Accept-Encoding", constant.AcceptEncoding)
	req.Header.Add("User-Agent", constant.UserAgent)
	req.Header.Add("Referer", constant.AcceptEncoding)
	req.Header.Add("Accept-Encoding", constant.AcceptEncoding)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("x-rpc-device_id", uuid.NewV4().String())
	//req.Body = ioutil.NopCloser(body)
	return
}

func addExtraGenshinHeader(add func(key, value string)) {
	add("x-rpc-client_type", constant.ClientType)
	add("x-rpc-app_version", constant.AppVersion)
	add("DS", util.GetDs())
}

func (g *GenshinClient) SendMessage(method, cookie, path, parameters string, extra bool, body, v interface{}) (err error) {
	var url string
	if parameters == "" {
		url = fmt.Sprintf("%s%s", constant.OpenApi, path)
	} else {
		url = fmt.Sprintf("%s%s%s", constant.OpenApi, path, parameters)
	}

	var jsonByte []byte = nil
	if body != nil {
		jsonByte, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("unable convert request body to json:%v", err)
		}
	}

	req, err := g.NewGenshinRequest(method, cookie, url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return err
	}
	if extra {
		addExtraGenshinHeader(req.Header.Add)
	}

	resp, err := g.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable convert response body io to json:%v", err)
	}
	log.Debug(string(respBody))

	if err = json.Unmarshal(respBody, &v); err != nil {
		return fmt.Errorf("unable convert response body json to struct:%v", err)
	}

	return
}

func (g *GenshinClient) SendGetMessage(cookie, path, parameters string, extra bool, v interface{}) (err error) {
	return g.SendMessage(http.MethodGet, cookie, path, parameters, extra, nil, v)
}

func (g *GenshinClient) SendPostMessage(cookie, path, parameters string, extra bool, body, v interface{}) (err error) {
	return g.SendMessage(http.MethodPost, cookie, path, parameters, extra, body, v)
}

/*
gResp *GenshinResponse, err error) {
	var url string
	if parameters == "" {
		url = fmt.Sprintf("%s%s", constant.OpenApi, path)
	} else {
		url = fmt.Sprintf("%s%s%s", constant.OpenApi, path, parameters)
	}

	var jsonByte []byte = nil
	if body != nil {
		jsonByte, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := g.NewGenshinRequest(http.MethodPost, cookie, url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return nil, err
	}
	if extra {
		addExtraGenshinHeader(req.Header.Add)
	}

	resp, err := g.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	log.Debug(string(b))
	if err != nil {
		return nil, err
	}

	return
*/
