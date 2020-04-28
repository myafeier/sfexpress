package sf_express

import (
	"fmt"

	"github.com/myafeier/log"
	"xorm.io/xorm"
)

var Daemon *Service

type Service struct {
	session *xorm.Session
}

func NewService(session *xorm.Session) *Service {
	return &Service{session: session}
}
func InitDaemon(session *xorm.Session, clientCode, clientCheckCode string) *Service {
	ClientCheckCode = clientCheckCode
	ClientCode = clientCode
	Daemon = &Service{
		session: session,
	}
	if has, err := session.IsTableExist(&SfExpressLog{}); err != nil {
		panic(err)
	} else {
		if !has {
			err := session.CreateTable(&SfExpressLog{})
			if err != nil {
				panic(err)
			}
			err = session.CreateIndexes(&SfExpressLog{})
			if err != nil {
				panic(err)
			}
		} else {
			err = session.Sync2(&SfExpressLog{})
			if err != nil {
				panic(err)
			}
		}
	}

	log.Info("sf express inited, clientCode:%s,clientCheckCode:%s", clientCode, clientCheckCode)
	return Daemon
}

func (s *Service) PostOrder(order *OrderRequestBody) (result *OrderResponseBody, err error) {
	sfLog := new(SfExpressLog)
	has, err := s.session.Where("inner_order_sn=?", order.OrderId).Get(sfLog)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if has {
		sfLog.Request = order
		_, err = s.session.ID(sfLog.Id).Cols("request").Update(sfLog)
		if err != nil {
			log.Error(err.Error())
			return
		}

	} else {
		sfLog.InnerOrderSn = order.OrderId
		sfLog.Request = order
		_, err = s.session.Insert(sfLog)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
	result, err = PostOrder(order)
	if err != nil || result.FilterResult != 2 {
		log.Error(err.Error())
		if result.FilterResult == 1 {
			err = fmt.Errorf("该订单需要顺丰人工确认是否可以邮寄")
		} else if result.FilterResult == 3 {
			err = fmt.Errorf("该订单无法邮寄，原因：%s", result.Remark)
		}
		sfLog.Response = new(OrderResponseBody)
		sfLog.Response.Remark = err.Error()
		s.session.ID(sfLog.Id).Cols("response").Update(sfLog)
		return
	}

	sfLog.Response = result
	_, err = s.session.ID(sfLog.Id).Cols("response").Update(sfLog)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}

func (s *Service) GetOne(outOrderSn string) (*SfExpressLog, error) {
	if outOrderSn == "" {
		err := fmt.Errorf("query param: outOrderSn is null")
		return nil, err
	}
	result := new(SfExpressLog)
	_, err := s.session.Where("inner_order_sn=?", outOrderSn).Get(result)
	return result, err
}

func (s *Service) GetRouteInfo(outOrderSn string) ([]RouteResponseBody, error) {
	req := new(RouteRequestBody)
	req.TrackingType = 2
	req.TrackingNumber = outOrderSn
	req.MethodType = 1
	return QueryRouteInfo(req)
}

func (s *Service) CheckOrderState(outOrderSn string) (*OrderResponseBody, error) {
	req := new(OrderStateRequestBody)
	req.OrderId = outOrderSn
	result, err := QueryOrderState(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	elog, err := s.GetOne(outOrderSn)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if elog.Id == 0 {
		err = fmt.Errorf("express order: %s  not found", outOrderSn)
		return nil, err
	}
	if elog.Response != nil && elog.Response.FilterResult != result.FilterResult {
		elog.Response.FilterResult = result.FilterResult
		elog.Response.Remark = result.Remark
		_, err = s.session.ID(elog.Id).Cols("response").Update(elog)
		if err != nil {
			log.Error(err.Error())
			return result, err
		}
	}
	return result, nil
}
