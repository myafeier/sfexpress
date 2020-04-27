package sf_express

import (
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
	if err != nil {
		log.Error(err.Error())
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
	result := new(SfExpressLog)
	_, err := s.session.Where("", outOrderSn).Get(result)
	return result, err
}
