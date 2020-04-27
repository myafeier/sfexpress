package sf_express

import (
	"encoding/xml"
	"testing"
)

func TestXMLEncode(t *testing.T) {

	req := new(OrderRequest)
	req.Service = "OrderService"
	req.Lang = "zh-CN"
	req.Head = "SLKJ2019"
	req.Body = OrderRequestBody{
		CustomerId: "xx",
		Cargo: Cargo{
			Name: "test",
		},
	}
	data, err := xml.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ecode:%s", data)

}

func TestResonseXML(t *testing.T) {
	res := new(OrderResponse)
	res.Service = "OrderService"
	res.Head = "xx"
	res.Body = OrderResponseBody{
		OrderId: "test",
	}

	data, err := xml.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ecode:%s", data)

}
