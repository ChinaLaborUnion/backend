package wx
//
//import (
//	"encoding/json"
//	constant "grpc-demo/constants/payment/wxpayment"
//	paymentEnums "grpc-demo/enums/payment"
//	"grpc-demo/models/db"
//	logUtils "grpc-demo/utils/log"
//
//	"github.com/iGoogle-ink/gopay/wechat"
//	"github.com/kataras/iris"
//)
//
//func RefundCallbackReceiver(ctx iris.Context) {
//	bs, err := wechat.ParseRefundNotify(ctx.Request())
//	if err != nil {
//		return
//	}
//	if bs.ReturnCode == wxpay.FAIL {
//		logUtils.Println(PaymentException.RefundCreateFail())
//		ctx.Text(replyMsg(wxpay.FAIL, ""))
//		return
//	}
//
//	requestJson, _ := json.Marshal(bs)
//	RawData := db.CallbackRawData{
//		Raw: string(requestJson),
//	}
//	db.Driver.Create(&RawData)
//
//	dcr, err := wechat.DecryptRefundNotifyReqInfo(bs.ReqInfo, constant.ApiKey)
//	if err != nil {
//		logUtils.Println(PaymentException.RefundCreateFail())
//		ctx.Text(replyMsg(wxpay.FAIL, "数据解密失败"))
//		return
//	} else {
//		logUtils.Println("数据解密:", dcr)
//	}
//
//	var payment db.WxPayOrder
//	getPayOrder := db.Driver.Where("out_trade_no = ?", dcr.OutTradeNo).First(&payment)
//	if getPayOrder != nil {
//		if payment.Status == paymentEnums.REFUND {
//			logUtils.Println(PaymentException.RefundCreateFail())
//			ctx.Text(replyMsg(wxpay.FAIL, "订单已处理"))
//			return
//		} else {
//			payment.Status = paymentEnums.REFUND
//
//		}
//		tx := db.Driver.Begin()
//		if err := tx.Debug().Save(&payment).Error; err != nil {
//			tx.Rollback()
//			logUtils.Println(PaymentException.PaymentSaveFail())
//			return
//		}
//		tx.Commit()
//	} else {
//		logUtils.Println(PaymentException.RefundCreateFail())
//		ctx.Text(replyMsg(wxpay.FAIL, ""))
//		return
//	}
//
//	var refund db.WxPayRefundOrder
//	getRefund := db.Driver.Where("order_id = ?", payment.OrderId).First(&refund)
//	if getRefund != nil {
//		status := dcr.RefundStatus
//		switch {
//		case status == "SUCCESS":
//			refund.Status = refundEnums.SUCCESS
//		case status == "CHANGE":
//			refund.Status = refundEnums.CHANGE
//		case status == "REFUNDCLOSE":
//			refund.Status = refundEnums.REFUNDCLOSE
//		}
//		refund.RefundStatus = dcr.RefundStatus
//		refund.RefundFee = gotil.String2Int(dcr.RefundFee)
//		refund.RefundAccount = dcr.RefundAccount
//		refund.RefundRecvAccout = dcr.RefundRecvAccout
//		refund.RefundRequestSource = dcr.RefundRequestSource
//		refund.SuccessTime = dcr.SuccessTime
//
//		tx := db.Driver.Begin()
//		if err := tx.Debug().Save(&refund).Error; err != nil {
//			tx.Rollback()
//			logUtils.Println(OrderException.OrderCreateFail())
//			return
//		}
//		tx.Commit()
//	} else {
//		logUtils.Println(PaymentException.QueryRefundFail())
//		return
//	}
//
//	var order db.TestOrder
//	//var childOrder []db.TestChildOrder
//
//	getOrder := db.Driver.Where("order_num = ?", payment.OutTradeNo).First(&order)
//	//getChildOrder := db.Driver.Where("order_id = ?", order.ID).Find(&childOrder)
//
//	if getOrder == nil {
//		logUtils.Println(OrderException.OrderIsNotExsit())
//		return
//	}
//
//	var orderExchange db.OrderExchange
//	if err := db.Driver.Where("order_id = ?", order.ID).First(&orderExchange).Error; err != nil {
//		panic(OrderException.ExchangeOrderIsNotExists())
//	} else {
//		orderExchange.Status = exchange.SUCCESS
//		tx := db.Driver.Begin()
//		if err := tx.Debug().Save(&orderExchange).Error; err != nil {
//			tx.Rollback()
//			logUtils.Println(PaymentException.PaymentSaveFail())
//			return
//		}
//		tx.Commit()
//	}
//
//	//状态更新
//	var childOrder db.TestChildOrder
//	if err := db.Driver.GetOne("test_child_order", orderExchange.ChildOrderId, &childOrder); err != nil {
//		logUtils.Println("无法找到子订单")
//		return
//	} //else {
//	//	childOrder.OrderStatus = orderEnums.AlreadyAS
//	//	tx := db.Driver.Begin()
//	//	if err := tx.Debug().Save(&childOrder).Error; err != nil {
//	//		tx.Rollback()
//	//		logUtils.Println(OrderException.OrderCreateFail())
//	//		return
//	//	}
//	//	tx.Commit()
//	//}
//
//	var orderDetail db.TestOrderDetail
//	if err := db.Driver.GetOne("test_order_detail", orderExchange.OrderDetailId, &orderDetail); err != nil {
//		logUtils.Println("无法找到订单详情")
//		return
//	} else {
//		var goods db.Goods
//		if err := db.Driver.GetOne("goods", orderDetail.GoodsID, &goods); err != nil {
//			logUtils.Println("无法找到商品")
//			return
//		} else {
//			goods.MonthSale -= orderDetail.PurchaseQty
//			goods.Total += orderDetail.PurchaseQty
//		}
//		tx := db.Driver.Begin()
//		if err := tx.Debug().Save(&goods).Error; err != nil {
//			tx.Rollback()
//			logUtils.Println(OrderException.OrderCreateFail())
//			return
//		}
//		tx.Commit()
//	}
//	orderDetail.IsPass = orderEnums.IsPassSUCCESS
//	tx := db.Driver.Begin()
//	if err := tx.Debug().Save(&orderDetail).Error; err != nil {
//		tx.Rollback()
//		logUtils.Println(OrderException.OrderCreateFail())
//		return
//	}
//	tx.Commit()
//
//	//tx := db.Driver.Begin()
//
//	//if err := tx.Debug().Save(&payment).Error; err != nil {
//	//	tx.Rollback()
//	//	logUtils.Println(PaymentException.PaymentSaveFail())
//	//	return
//	//}
//
//	ctx.Text(replyMsg(wxpay.SUCCESS, "OK"))
//}