package sf_express

import (
	"encoding/xml"
	"time"
)

const (
	Url                   = "http://bsp-oisp.sf-express.com/bsp-oisp/sfexpressService" //请求地址
	InterfaceOfOrder      = "OrderService"                                             //下单接口名称
	InterfaceOfOrderState = "OrderSearchService"                                       //订单结果查询接口
	InterfaceOfRouteInfo  = "RouteService"                                             //路由查询接口
)
const (
	OrderFilterResultOfNeedVerify OrderFilterResult = 1
	OrderFilterResultOfOk         OrderFilterResult = 2
	OrderFilterResultOfDeny       OrderFilterResult = 3
)

type OrderFilterResult int8

var OrderFilterResults = map[OrderFilterResult]string{
	OrderFilterResultOfNeedVerify: "待人工确认",
	OrderFilterResultOfOk:         "可收派",
	OrderFilterResultOfDeny:       "不可收派",
}

func (o OrderFilterResult) ToString() string {
	if r, ok := OrderFilterResults[o]; ok {
		return r
	} else {
		return "-"
	}
}

//筛单结果:1:人工确认2:可收派3:不可以收派
var ClientCode = ""      //客户编码
var ClientCheckCode = "" //客户校验码

//下单请求
type OrderRequest struct {
	Request
	Body OrderRequestBody `xml:"Body>Order"`
}

type OrderResponse struct {
	Response
	Body OrderResponseBody `xml:"Body>OrderResponse"`
}

// 订单状态查询响应
type OrderStateResponse struct {
	Response
	Body OrderResponseBody `xml:"Body>OrderResponse"`
}

// 订单状态查询请求
type OrderStateRequest struct {
	Request
	Body OrderStateRequestBody `xml:"Body>OrderSearch"`
}

// 路由查询
type RouteRequest struct {
	Request
	Body RouteRequestBody `xml:"Body>RouteRequest"`
}

//路由查询响应
type RouteResponse struct {
	Response
	Body []RouteResponseBody `xml:"Body>RouteResponse"`
}

// 下单请求
type OrderRequestBody struct {
	XMLName           xml.Name `xml:"Order"`
	OrderId           string   `xml:"orderid,attr"`                 //必填,客户订单号,建议英文字母+YYMMDD(日期)+流水号,如:TB1207300000001
	FromCompany       string   `xml:"j_company,attr"`               //寄件方公司名称,如果需要生成电子面单,则为必填。
	FromContact       string   `xml:"j_contact,attr"`               // 寄件方联系人,如果需要生成电子面单,则为必填。
	FromPhone         string   `xml:"j_tel,attr"`                   //寄件方联系电话,如果需要生成电子面单,则为必填。
	FromProvince      string   `xml:"j_province,attr"`              //寄件方所在省份字段填写要求:必须是标准的省名称称谓 如:广东省,如果是直辖市,请直接传北京、上海等。
	FromCity          string   `xml:"j_city,attr"`                  //寄件方所在城市名称,字段填写k要求:必须是标准的城市称谓如:深圳市。
	FromCounty        string   `xml:"j_county,attr"`                //寄件人所在县/区,必须是标准的县/区称谓,示例:“福田区”。
	FromAddress       string   `xml:"j_address,attr"`               //	寄件方详细地址,包括省市区,示例:“广东省深圳市福田区新洲十一街万基商务大厦10楼” ,如果需要生成电子面单,则必填。
	ToCompany         string   `xml:"d_company,attr"`               //
	ToContact         string   `xml:"d_contact,attr"`               //
	ToTel             string   `xml:"d_tel,attr"`                   //
	ToProvince        string   `xml:"d_province,attr"`              //
	ToCity            string   `xml:"d_city,attr"`                  //
	ToCounty          string   `xml:"d_county,attr"`                //
	ToAddress         string   `xml:"d_address,attr"`               //
	CustomerId        string   `xml:"custid,attr"`                  //顺丰月结卡号
	ExpressType       string   `xml:"express_type,attr"`            // 快件产品编码,详见附录 11	T4	医药安心递
	SendStartTime     string   `xml:"sendstarttime,attr,omitempty"` // 要求上门取件开始时间,格式:YYYY-MM-DD HH24:MM:SS 示例:2012-7-30 09:30:00。
	DoCall            int8     `xml:"is_docall,attr"`               //是否要求通过手持终端通知顺丰收派员收件:1:要求其它为不要求
	UseUnifiedWayBill int8     `xml:"is_unified_waybill_no,attr"`   //是否使用国家统一面单号:默认0 。1:是, 0:否
	Cargo             Cargo    `xml:"Cargo"`                        //货物信息
}

//货物信息
type Cargo struct {
	Name  string `xml:"name,attr"`            //货物名称,如果需要生成电子面单,则为必填。
	Count int8   `xml:"count,attr,omitempty"` // 货物数量跨境件报关需要填写
}

//下单响应
type OrderResponseBody struct {
	XMLName    xml.Name `xml:"OrderResponse"`
	OrderId    string   `xml:"orderid,attr"`    //必填,客户订单号,建议英文字母+YYMMDD(日期)+流水号,如:TB1207300000001
	MailNo     string   `xml:"mailno,attr"`     //顺丰运单号,一个订单只能有一个母单号,如果是子母单的情况,以半角逗号分隔,主单号在第一个位置
	OriginCode string   `xml:"origincode,attr"` //原寄地区域代码
	DestCode   string   `xml:"destcode,attr"`   //目的地区域代码

	FilterResult OrderFilterResult `xml:"filter_result,attr"` //筛单结果:1:人工确认2:可收派3:不可以收派
	Remark       string            `xml:"remark,attr"`        //	filter_result=3时必填,不可以收派的原因代码:1:收方超范围2:派方超范围3:其它原因高峰管控提示信息【数字】:【高峰管控提示信息】(如 4:温馨提示 ,1:春运延时)
}

// 查询订单的处理结果
type OrderStateRequestBody struct {
	XMLName xml.Name `xml:"OrderSearch"`
	OrderId string   `xml:"orderid,attr"`
}

//路由查询请求
type RouteRequestBody struct {
	XMLName        xml.Name `xml:"RouteRequest"`
	TrackingType   int8     `xml:"tracking_type,attr"`   //2:根据客户订单号查询,order节点中tracking_number将被当作客户订单号处理
	TrackingNumber string   `xml:"tracking_number,attr"` //如果tracking_type=1,则此值为顺丰运单号,如果有多个单号,以逗号分隔,如"123,124,125"。
	MethodType     int8     `xml:"method_type,attr"`     //1:标准路由查询
}

//路由查询响应
type RouteResponseBody struct {
	XMLName    xml.Name    `xml:"RouteResponse"`
	OrderId    string      `xml:"orderid,attr"` //必填,客户订单号,建议英文字母+YYMMDD(日期)+流水号,如:TB1207300000001
	MailNo     string      `xml:"mailno,attr"`  //顺丰运单号,一个订单只能有一个母单号,如果是子母单的情况,以半角逗号分隔,主单号在第一个位置
	RouteInfos []RouteInfo `xml:"Route"`
}
type RouteInfo struct {
	XMLName       xml.Name `xml:"Route"`
	AcceptTime    string   `xml:"accept_time,attr" json:"accept_time"`       //路由节点发生的时间,格式:YYYY-MM-DD HH24:MM:SS,示例:2012-7-30 09:30:00
	AcceptAddress string   `xml:"accept_address,attr" json:"accept_address"` //路由节点发生的地点
	Remark        string   `xml:"remark,attr" json:"remark"`                 //路由节点具体描述
	Opcode        string   `xml:"opcode,attr" json:"opcode"`                 //路由节点操作码
}

//通用查询结构
type Request struct {
	XMLName xml.Name `xml:"Request"`
	Service string   `xml:"service,attr"`
	Lang    string   `xml:"lang,attr"`
	Head    string   `xml:"Head"`
}

// 通用响应结构
type Response struct {
	XMLName xml.Name `xml:"Response"`
	Service string   `xml:"service,attr"`
	Error   struct {
		Code string `xml:"code,attr"`
	} `xml:"ERROR,omitempty"`
	Head string `xml:"Head"`
}

//顺丰日志数据库
type SfExpressLog struct {
	Id           int64              `json:"id"`
	InnerOrderSn string             `json:"inner_order_sn" xorm:"varchar(100) default '' index"` //内部订单号
	SfOrderSn    string             `json:"sf_order_sn" xorm:"varchar(100) default '' index"`    //顺丰订单号
	Request      *OrderRequestBody  `json:"request" xorm:"json"`                                 //请求
	Response     *OrderResponseBody `json:"response" xorm:"json"`                                //响应
	Created      time.Time          `json:"created" xorm:"created"`
	Updated      time.Time          `json:"updated" xorm:"updated"`
}

func (self *SfExpressLog) TableName() string {
	return "sf_express_log"
}
