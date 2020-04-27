package sf_express

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/myafeier/log"
)

func QueryRouteInfo(req *RouteRequestBody) (result []RouteResponseBody, err error) {
	TData := new(RouteRequest)
	TData.Request = Request{}
	TData.Service = InterfaceOfRouteInfo
	TData.Head = ClientCode
	TData.Lang = "zh-CN"
	TData.Body = *req

	data, err := xml.MarshalIndent(TData, "", "")
	if err != nil {
		log.Error(err.Error())
		return
	}

	param := url.Values{}
	param.Set("xml", fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8'?>%s`, data))
	param.Set("verifyCode", signature(data))

	log.Debug("param xml:%s, verifyCode:%s", data, signature(data))

	resp, err := http.PostForm(Url, param)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code :%d", resp.StatusCode)
		return
	}

	var OrginData RouteResponse
	err = xml.Unmarshal(data, &OrginData)
	if err != nil {
		log.Error("error:%s,data:%s", err.Error(), data)
		return

	}
	log.Debug("data: %s\n stru: %+v", data, OrginData)
	if OrginData.Head != "OK" {
		err = fmt.Errorf("下单失败，顺丰返回:%s", OrginData.Error.Code)
		return
	}

	result = OrginData.Body
	return
}
