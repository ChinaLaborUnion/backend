package transaction

import (
	"bytes"
	authbase "grpc-demo/core/auth"
	"grpc-demo/exceptions/goods"
	"grpc-demo/exceptions/order"
	"grpc-demo/models/db"
	util "grpc-demo/utils/hash"
	"grpc-demo/views/transaction/payment/ali"
	"grpc-demo/views/transaction/payment/wx"

	"github.com/kataras/iris"
)

func PaymentMiddleware(ctx iris.Context, auth authbase.AuthAuthorization, oid, tid int) {
	//auth.CheckLogin()
	var order db.OrderInfo
	if err := db.Driver.GetOne("order", oid, &order); err != nil {
		panic(orderException.OrderNotExist())
	}
	//if order.AccountID != auth.AccountModel().Id || auth.IsAdmin() {
	//	panic(accountException.NoPermission())
	//}


	//读取商品信息
	//使用bytesBuffer拼接字符串
	var buf = bytes.Buffer{}

	var goods db.GoodsInfo
	if err := db.Driver.GetOne("goods_info", order.GoodsID, &goods); err != nil {
		panic(goodsException.GoodsNotExist())
	} else {
		buf.WriteString("["+goods.Name+"]"+"/"+goods.Brief)
	}



	OutTradeNo := order.Number
	GoodsInfo := buf
	AccountId := auth.AccountModel().Id
	aliTotalAmount := util.Int2String(order.TotalPrice)
	wxTotalAmount := order.TotalPrice
	//openId := auth.AccountModel().OpenId

	switch tid {
	case 1:
		ali.PaymentForAli(ctx, oid, AccountId, OutTradeNo, GoodsInfo.String(), aliTotalAmount)

	case 2:
		wx.PaymentForWx(ctx, oid, AccountId, wxTotalAmount, OutTradeNo, GoodsInfo.String())
	}

}

func RefundMiddleware(ctx iris.Context){

}
