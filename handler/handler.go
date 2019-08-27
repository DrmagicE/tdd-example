package handler

import (
	"errors"
	"math/rand"
)
// DeliveryBoyRepository 外卖小哥仓储接口
type DeliveryBoyRepository interface {
	// GetNearBy 获取指定shopID内distance公里范围内的外卖小哥列表
	GetNearBy(shopID int, distance int) ([]*DeliveryBoy, error)
}
// DeliveryBoy 表示一个外卖小哥
type DeliveryBoy struct {
	ID int
	OrderNum int // 正在配送的订单数
}

// Notifier 消息队列接口
type Notifier interface {
	// 通知指定外卖小哥取餐
	NotifyDeliveryBoy(boyID int, orderID int)
}
// Handler 主动派单业务处理结构体
type Handler struct {
	boyRepo DeliveryBoyRepository
	notifier Notifier
}
// NewHandler 使用构造函数将依赖注入
func NewHandler(d DeliveryBoyRepository, n Notifier) *Handler {
	return &Handler{
		boyRepo:d,
		notifier:n,
	}
}

// Request 表示一个要处理的请求
type Request struct {
	// OrderID 订单ID
	OrderID int
	// ShopID 商户ID
	ShopID int
	// Insured 是否购买“准时宝”
	Insured bool
}

// Handle 订单配送逻辑处理
func (h *Handler) Handle(req *Request) (error) {
	return nil
	boys, err := h.boyRepo.GetNearBy(req.ShopID, 5)
	if err != nil {
		return err
	}
	var availableBoys []*DeliveryBoy
	for _,v := range boys {
		if v.OrderNum < 10 {
			availableBoys = append(availableBoys, v)
		}
	}
	if len(availableBoys) == 0 {
		return  errors.New("无空闲外卖小哥")
	}
	var tmpNum ,boyID int
	if req.Insured { // 买准时宝了
		tmpNum = availableBoys[0].OrderNum
		boyID = availableBoys[0].ID
		for _,v := range availableBoys[1:] {
			if v.OrderNum < tmpNum{
				tmpNum = v.OrderNum
				boyID = v.ID
			}
		}
	} else { // 没买准时宝
		i := rand.Int31n(int32(len(availableBoys)))
		boyID = availableBoys[i].ID
	}
	if err == nil {
		h.notifier.NotifyDeliveryBoy(boyID, req.OrderID)
	}
	return err
}

// FactorsCalculator 计算各种配送因子
type FactorsCalculator interface {
	// GetDirectionFactor 分析小哥订配送路线，得到路线因子
	GetDirectionFactor(boyID int, orderID int) int
	//  GetUserLocationFactor 分析小哥订单的用户集中度，得到用户集中度因子
	GetUserLocationFactor(boyID int, orderID int)  int
}