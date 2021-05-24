package ali

import (
	constant "grpc-demo/constants/payment/alipayment"
	transactionEnums "grpc-demo/enums/transaction"
	"grpc-demo/exceptions/db"

	"grpc-demo/enums/order"
	paymentEnums "grpc-demo/enums/payment"
	"grpc-demo/exceptions/order"
	"grpc-demo/exceptions/transaction"
	util "grpc-demo/utils/hash"

	"grpc-demo/models/db"
	logUtils "grpc-demo/utils/log"

	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/alipay"
	"github.com/kataras/iris"
)

func PaymentCallbackReceiver(ctx iris.Context) {
	bm := make(gopay.BodyMap)
	bm, err := alipay.ParseNotifyToBodyMap(ctx.Request())
	if err != nil {
		ctx.Write([]byte("fail"))
		return
	} else {
		//验证签名
		ok, err := alipay.VerifySign(constant.AliPayPublicKey, bm)
		if err != nil {
			logUtils.Println("err:", err)
		}
		if ok != true {
			ctx.Write([]byte("fail"))
			return
		}
		logUtils.Println("异步验签是否通过:", ok)
	}

	if err := bm.Get("trade_status"); err != "TRADE_SUCCESS" {
		logUtils.Println("微信通知接口接收通信失败,交易状态说明：%s", err)
		ctx.Write([]byte("fail"))

	} else {
		var order db.OrderInfo
		if err := db.Driver.Where("number = ?", bm.Get("out_trade_no")).First(&order).Error; err == nil {
			logUtils.Println(orderException.OrderNotExist())
			return
		}

		var payment db.AliPayOrder
		if err := db.Driver.Where("out_trade_no = ?", bm.Get("out_trade_no")).First(&payment).Error; err != nil {
			logUtils.Println(transactionException.PaymentException())
			return
		}
		if state := payment.TradeStatus; state == "TRADE_SUCCESS" {
			logUtils.Println("回调已处理")
			return
		}
		if totalFee := bm.Get("total_amount"); totalFee != payment.TotalAmount && totalFee != util.Int2String(order.TotalPrice) {
			logUtils.Println(transactionException.AmoutnIsNotEqual())
			ctx.Write([]byte("fail"))
			return
		}

		tx := db.Driver.Begin()
		order.Status = orderEnum.Done
		if err := tx.Debug().Save(&order).Error; err != nil {
			tx.Rollback()
			logUtils.Println(dbException.SaveFail())
			return
		}

		payment.Status = paymentEnums.SUCCESS

		if err := tx.Save(&payment).Error; err != nil {
			tx.Rollback()
			logUtils.Println(transactionException.PaymentSaveFail())
			return
		}

		var transaction db.TransactionInfo
		if err := db.Driver.Where("order_id = ?", order.ID).First(&transaction).Error;err!=nil {
			logUtils.Println(transactionException.TransactionGetFail())
		}else{
			transaction.TransactionStatus = transactionEnums.PAID
			if err := tx.Save(&payment).Error; err != nil {
				tx.Rollback()
				logUtils.Println(transactionException.PaymentSaveFail())
				return
			}
		}

		tx.Commit()

	}
	ctx.Write([]byte("success"))

}