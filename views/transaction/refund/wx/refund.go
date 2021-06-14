package wx

import (
	authbase "grpc-demo/core/auth"
	PaymentEnums "grpc-demo/enums/payment"
	accountException "grpc-demo/exceptions/account"
	"grpc-demo/models/db"
	logUtils "grpc-demo/utils/log"

	"github.com/iGoogle-ink/gopay"
	"github.com/kataras/iris"
)

func RefundForWx(ctx iris.Context, auth authbase.AuthAuthorization, eid int) {
	auth.CheckLogin()
	if auth.IsAdmin() != true {
		panic(accountException.NoPermission())
	}
	var orderExchange db.OrderExchange
	if err := db.Driver.GetOne("order_exchange", eid, &orderExchange); err != nil {
		panic(OrderException.OrderIsNotExsit())
	}
	if orderExchange.Status != exchange.PERMIT {
		panic(PaymentException.StatusIsNotAllow())
	}

	var payment db.WxPayOrder
	if err := db.Driver.Where("order_id = ?", orderExchange.OrderId).First(&payment).Error; err != nil {
		panic(PaymentException.CanNotFind())
	} else {
		if payment.Status == PaymentEnums.CANCEL {
			panic(PaymentException.PaymentException())
		}
	}

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", payment.OutTradeNo)
	bm.Set("nonce_str", hash.GetRandomString(16))
	bm.Set("sign_type", wxpay.SignType_HMAC_SHA256)
	bm.Set("out_refund_no", orderExchange.OutRefundNo)
	bm.Set("total_fee", payment.TotalFee)
	bm.Set("refund_fee", orderExchange.ReturnAmount)
	bm.Set("notify_url", constant.RefundNotifyUrl)
	bm.Set("refund_account", "REFUND_SOURCE_RECHARGE_FUNDS")

	//请求申请退款（沙箱环境下，证书路径参数可传空）
	//    body：参数Body
	//    certFilePath：cert证书路径
	//    keyFilePath：Key证书路径
	//    pkcs12FilePath：p12证书路径
	wxRsp, resBm, err := client.Refund(bm, nil, nil, nil)
	if err != nil {
		logUtils.Println("订单错误:", err)
		return
	}
	if resBm != nil {
		logUtils.Println("订单数据:", resBm)
	}

	var refund = db.WxPayRefundOrder{
		AccountId:     orderExchange.AccountId,
		OrderId:       orderExchange.OrderId,
		TransactionId: payment.TransactionId,
		OutTradeNo:    payment.OutTradeNo,
		OutRefundNo:   orderExchange.OutRefundNo,
		TotalFee:      payment.TotalFee,
		Status:        RefundEnums.PROCESSING,
		RefundId:      wxRsp.RefundId,
	}
	payment.Status = PaymentEnums.REFUNDING
	tx := db.Driver.Begin()
	if err := tx.Debug().Create(&refund).Error; err != nil {
		tx.Rollback()
		panic(PaymentException.RefundCreateFail())
	}
	if err := tx.Debug().Save(&payment).Error; err != nil {
		tx.Rollback()
		panic(PaymentException.RefundCreateFail())
	}

	tx.Commit()
	ctx.JSON(iris.Map{
		"respond": iris.Map{
			"ResultCode": wxRsp.ResultCode,
			"Msg":        wxRsp.ReturnMsg,
			"ReturnCode": wxRsp.ReturnCode,
			"ErrMessage": wxRsp.ErrCodeDes,
		},
	})

}