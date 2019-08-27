package handler

import (
	"testing"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)


// 1. 商户5公里内没有外卖小哥存在时，返回错误，不执行后续派单操作
func TestHandler_Handle_NoBoy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer  ctrl.Finish()
	a := assert.New(t)

	// 使用gomock 来mock依赖
	d := NewMockDeliveryBoyRepository(ctrl)
	n := NewMockNotifier(ctrl)
	h := NewHandler(d,n)

	req := &Request{
		OrderID:1,
		ShopID:2,
	}
	// 5公里内没有外卖小哥
	d.EXPECT().GetNearBy(req.ShopID, 5).Return(nil, errors.New("o no..5公里内没有外卖小哥"))
	err := h.Handle(req)
	a.Error(err)
}

// 2. 商户5公里内所有的外卖小哥配送订单数均>=10时，返回错误，不执行后续派单操作。
func TestHandler_Handle_NoAvailableBoy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer  ctrl.Finish()
	a := assert.New(t)

	d := NewMockDeliveryBoyRepository(ctrl)
	n := NewMockNotifier(ctrl)
	h := NewHandler(d,n)
	req := &Request{
		OrderID:1,
		ShopID:2,
	}
	// 返回的外卖小哥配送订单均>=10
	d.EXPECT().GetNearBy(req.ShopID, 5).Return([]*DeliveryBoy{
		{
			ID:1,
			OrderNum:10,
		},
		{
			ID:2,
			OrderNum:11,
		},

	}, nil)
	err := h.Handle(req)
	a.Error(err)
}

// 3. 如下单用户购买了准时宝，选择一个订单最少的小哥，通知小哥取餐。
func TestHandler_Handle_Insured(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer  ctrl.Finish()
	a := assert.New(t)

	d := NewMockDeliveryBoyRepository(ctrl)
	n := NewMockNotifier(ctrl)

	h := NewHandler(d,n)
	req := &Request{
		OrderID:1,
		ShopID:2,
		Insured:true,
	}
	d.EXPECT().GetNearBy(req.ShopID, 5).Return([]*DeliveryBoy{
		{
			ID:1,
			OrderNum:4,  // <-- 这个应当是被选中的外卖小哥
		},
		{
			ID:2,
			OrderNum:5,
		},
		{
			ID:3,
			OrderNum:6,
		},

	}, nil)

	// 通知小哥取餐
	n.EXPECT().NotifyDeliveryBoy(1, req.OrderID)
	err := h.Handle(req)
	a.Nil(err)
}

// 4. 如下单用户未购买准时宝，随机分配一名小哥服务。通知小哥取餐，商铺备餐。
func TestHandler_Handle_NotInsured(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer  ctrl.Finish()
	a := assert.New(t)

	d := NewMockDeliveryBoyRepository(ctrl)
	n := NewMockNotifier(ctrl)

	h := NewHandler(d,n)
	req := &Request{
		OrderID:1,
		ShopID:2,
		Insured:true,
	}
	d.EXPECT().GetNearBy(req.ShopID, 5).Return([]*DeliveryBoy{
		{
			ID:1,
			OrderNum:4,
		},
		{
			ID:2,
			OrderNum:5,
		},
		{
			ID:3,
			OrderNum:6,
		},

	}, nil)
	orderID := 1
	// 通知小哥取餐
	n.EXPECT().NotifyDeliveryBoy(gomock.Any(), orderID)
	err := h.Handle(req)
	a.Nil(err)
}