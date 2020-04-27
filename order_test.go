package sf_express

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/myafeier/log"
)

func init() {
	log.SetLogLevel(log.DEBUG)
	ClientCheckCode = "***"
	ClientCode = "***"
}
func TestOrder(t *testing.T) {

	order := OrderRequestBody{}
	order.CustomerId = "7551234567"
	order.DoCall = 1
	order.ExpressType = "T4"
	order.FromAddress = "江东华城"
	order.FromProvince = "云南省"
	order.FromCity = "昆明市"
	order.FromCounty = "盘龙区"
	order.FromCompany = "测试公司"
	order.FromContact = "夏菲"
	order.FromPhone = "18987092111"
	order.OrderId = fmt.Sprintf("To%d", rand.Int63n(9999999999999999))
	order.SendStartTime = "2020-04-01 15:33:33"
	order.ToAddress = "海龟创业园"
	order.ToProvince = "云南省"
	order.ToCity = "昆明市"
	order.ToCounty = "呈贡区"
	order.ToCompany = "王三"
	order.ToContact = "王4"
	order.ToTel = "18987092884"
	order.UseUnifiedWayBill = 1
	order.Cargo = Cargo{
		Name: "测试商品",
	}
	result, err := PostOrder(&order)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resp: %+v", *result)

}
func TestQueryOrderState(t *testing.T) {
	req := new(OrderStateRequestBody)
	req.OrderId = "1000100000000172"
	resp, err := QueryOrderState(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", *resp)

}
