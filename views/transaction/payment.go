package transaction

import (
	"bytes"
	authbase "grpc-demo/core/auth"
	accountException "grpc-demo/exceptions/account"
	"grpc-demo/exceptions/goods"
	"grpc-demo/exceptions/order"
	transactionException "grpc-demo/exceptions/transaction"
	"grpc-demo/models/db"
	util "grpc-demo/utils/hash"
	logUtils "grpc-demo/utils/log"
	"grpc-demo/views/transaction/payment/ali"
	"grpc-demo/views/transaction/payment/wx"
	"grpc-demo/views/transaction/payment/wxv2"

	"github.com/kataras/iris"
)

func PaymentMiddleware(ctx iris.Context, auth authbase.AuthAuthorization, oid, tid int) {
	auth.CheckLogin()
	var order db.OrderInfo
	if err := db.Driver.GetOne("order", oid, &order); err != nil {
		panic(orderException.OrderNotExsit())
	}
	if order.AccountID != auth.AccountModel().Id || auth.IsAdmin() {
		panic(accountException.NoPermission())
	}


	//读取商品信息
	//使用bytesBuffer拼接字符串
	var buf = bytes.Buffer{}

	var goods db.GoodsInfo
	if err := db.Driver.GetOne("goods_info", order.GoodsID, &goods); err != nil {
		panic(goodsException.GoodsNotExsit())
	} else {
		buf.WriteString("["+goods.Name+"]")
	}



	OutTradeNo := order.Number
	GoodsInfo := buf
	AccountId := auth.AccountModel().Id
	if order.TotalPrice<0{
		panic(transactionException.TransactionCreateFail())
	}
	aliTotalAmount := util.Float32ToString(util.Save2Decimal(float64(order.TotalPrice)))
	wxTotalAmount := order.TotalPrice
	//openId := auth.AccountModel().OpenId

	logUtils.Println(aliTotalAmount)

	switch tid {
	case 1:
		ali.PaymentForAli(ctx, oid, AccountId, OutTradeNo, GoodsInfo.String(), aliTotalAmount)

	case 2:
		wx.PaymentForWx(ctx, oid, AccountId, wxTotalAmount, OutTradeNo, GoodsInfo.String())

	case 3:
		wxv2.PaymentForWxV2(ctx, oid, AccountId, wxTotalAmount, OutTradeNo, GoodsInfo.String())

	default:
		panic(transactionException.StatusIsNotAllow())
	}

}


