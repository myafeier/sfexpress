package sf_express

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/myafeier/log"
)

func PostOrder(order *OrderRequestBody) (result *OrderResponseBody, err error) {
	TData := new(OrderRequest)
	TData.Service = InterfaceOfOrder
	TData.Head = ClientCode
	TData.Lang = "zh-CN"
	TData.Body = *order

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

	var OrginData OrderResponse
	err = xml.Unmarshal(data, &OrginData)
	if err != nil {
		log.Error("error:%s,data:%s", err.Error(), data)
		return

	}
	log.Debug("%s", data)
	if OrginData.Head != "OK" {
		err = fmt.Errorf("下单失败，顺丰返回:%s", OrginData.Error.Code)
		return
	}
	result = &OrginData.Body
	return
}

func QueryOrderState(req *OrderStateRequestBody) (result *OrderResponseBody, err error) {
	TData := new(OrderStateRequest)
	TData.Request = Request{}
	TData.Service = InterfaceOfOrderState
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

	var OrginData *OrderResponse
	err = xml.Unmarshal(data, &OrginData)
	if err != nil {
		log.Error("error:%s,data:%s", err.Error(), data)
		return

	}
	log.Debug("%s", data)
	if OrginData.Head != "OK" {
		err = fmt.Errorf("查询失败，顺丰返回:%s", OrginData.Error.Code)
		return
	}

	result = &OrginData.Body
	return
}
func signature(data []byte) string {
	md5sign := md5.New()
	md5sign.Write([]byte(fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8'?>%s%s`, data, ClientCheckCode)))
	return base64.StdEncoding.EncodeToString(md5sign.Sum(nil))
}
