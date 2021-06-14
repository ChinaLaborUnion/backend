package wx

import (
	"encoding/json"
	"fmt"
	constant "grpc-demo/constants/payment/wxpayment"
	wxconstant "grpc-demo/constants/payment/wxpayment"
	orderEnum "grpc-demo/enums/order"
	PaymentEnums "grpc-demo/enums/payment"
	transactionEnums "grpc-demo/enums/transaction"
	orderException "grpc-demo/exceptions/order"
	"grpc-demo/exceptions/transaction"
	"grpc-demo/models/db"
	util "grpc-demo/utils/hash"
	logUtils "grpc-demo/utils/log"
	"grpc-demo/views/transaction/payutil"

	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/kataras/iris"
)

type PayerStruct struct {
	Openid             string `xml:"openid" json:"openid"`
}
type AmountStruct struct {
	Total             int    `xml:"total" json:"total"`
	PayerTotal        int    `xml:"payer_total" json:"payer_total"`
	Currency          string `xml:"currency" json:"currency"`
	PayerCurrency     string `xml:"payer_currency" json:"payer_currency"`
}
type SceneInfoStruct struct {
	DeviceId           string `xml:"device_id" json:"device_id"`
}
type PromotionDetailStruct struct {
	CouponId           string `xml:"coupon_id" json:"coupon_id"`
	Name               string `xml:"name" json:"name"`
	Scope              string `xml:"scope" json:"scope"`
	type1              string `xml:"type"  json:"type"`
	Amount             int  `xml:"amount" json:"amount"`
	StockId            string `xml:"stock_id" json:"stock_id"`
	WechatpayContribute string `xml:"wechatpay_contribute" json:"wechatpay_contribute"`
	MerchantContribute string  `xml:"merchant_contribute" json:"merchant_contribute"`
	OtherContribute    string  `xml:"other_contribute" json:"other_contribute"`
	Currency           string `xml:"currency" json:"currency"`
}
type GoodsDetailSturct struct {
	GoodId             string `xml:"good_id" json:"good_id"`
	Quantity           int    `xml:"quantity" json:"quantity"`
	UnitPrice          int    `xml:"unit_price" json:"unit_price"`
	DiscountAmount     int    `xml:"discount_amount" json:"discount_amount"`
	GoodsRemark        string `xml:"goods_remark" json:"goods_remark"`
}
type Resource struct {
	Appid              string `xml:"appid,omitempty" json:"appid,omitempty"`
	MchId              string `xml:"mch_id,omitempty" json:"mch_id,omitempty"`
	OutTradeNo         string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty"`
	TransactionId      string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	TradeType          string `xml:"trade_type,omitempty" json:"trade_type,omitempty"`
	TradeState		   string `xml:"trade_state" json:"trade_state"`
	TradeStateDesc     string `xml:"trade_state_desc" json:"trade_state_desc"`
	BankType           string `xml:"bank_type,omitempty" json:"bank_type,omitempty"`
	Attach             string `xml:"attach,omitempty" json:"attach,omitempty"`
	SuccessTime        string `xml:"success_time" json:"success_time"`
	SubMchId           string `xml:"sub_mch_id,omitempty" json:"sub_mch_id,omitempty"`
	DeviceInfo         string `xml:"device_info,omitempty" json:"device_info,omitempty"`
	NonceStr           string `xml:"nonce_str,omitempty" json:"nonce_str,omitempty"`
	Sign               string `xml:"sign,omitempty" json:"sign,omitempty"`
	SignType           string `xml:"sign_type,omitempty" json:"sign_type,omitempty"`
	Openid             string `xml:"openid,omitempty" json:"openid,omitempty"`
	IsSubscribe        string `xml:"is_subscribe,omitempty" json:"is_subscribe,omitempty"`
	Payer              PayerStruct 	`xml:"payer" json:"payer"`
	Amount             AmountStruct `xml:"amount" json:"amount"`
	SceneInfo          SceneInfoStruct `xml:"scene_info" json:"scene_info"`
	PromotionDetail    PromotionDetailStruct `xml:"promotion_detail" json:"promotion_detail"`
}

func CallbackReceiver(ctx iris.Context) {
	request,err := payutil.ParseNotify(ctx.Request())
	//request, err := wechat.ParseNotify(ctx.Request())
	requestJson, _ := json.Marshal(request)
	logUtils.Println(requestJson)
	if err != nil {
		ctx.Text(replyMsg("FAIL", "解析错误"))
		return
	}
	decryptBytes, err := payutil.DecryptGCM(
		wxconstant.AesKey,
		request.Resource.Nonce,
		request.Resource.Ciphertext,
		request.Resource.AssociatedData)
	if err!= nil{
		ctx.Text(replyMsg("FAIL", "解析错误"))
		return
	}

	resource := new(Resource)
	if err := json.Unmarshal(decryptBytes, resource);err!=nil{
		ctx.Text(replyMsg("FAIL", "解析错误"))
	}
	logUtils.Println(resource)

	if err := resource.TradeState; err != "SUCCESS" {
		logUtils.Println("微信通知接口接收通信失败:", resource.TradeStateDesc)
		ctx.Text(replyMsg("FAIL", "状态不符"))

	} else {


		var order db.OrderInfo
		if err := db.Driver.Where("number = ?", resource.OutTradeNo).First(&order).Error;err==nil{
			logUtils.Println(orderException.OrderNotExsit())
			return
		}

		var payment db.WxPayOrder
		if err := db.Driver.Where("out_trade_no = ?", resource.OutTradeNo).First(&payment).Error;err!=nil{
			payment.TradeStateDesc = resource.TradeStateDesc
			logUtils.Println(transactionException.PaymentException())
			return
		}

		if state := payment.TradeState; state == "SUCCESS" {
			logUtils.Println("回调已处理")
			return
		}
		if totalFee := resource.Amount.Total; totalFee != payment.TotalFee && totalFee != order.Total {
			logUtils.Println(transactionException.AmoutnIsNotEqual())
			ctx.Text(replyMsg("FAIL", "金额不等"))
			return
		}

		tx := db.Driver.Begin()
		order.Status = orderEnum.Done
		if err := tx.Debug().Save(&order).Error; err != nil {
			tx.Rollback()
			logUtils.Println(orderException.CancelRefuse())
			return
		}

		payment.Status = PaymentEnums.SUCCESS
		payment.TradeState = resource.TradeState
		payment.TradeStateDesc = resource.TradeStateDesc
		payment.TransactionId = resource.TransactionId
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
	ctx.Text(replyMsg("SUCCESS", "OK"))
}

func PaymentCallbackReceiver(ctx iris.Context) {
	request, err := wechat.ParseNotify(ctx.Request())
	requestJson, _ := json.Marshal(request)

	logUtils.Println(string(requestJson))

	if err != nil {
		ctx.Text(replyMsg("FAIL", "解析错误"))

		return
	} else {
		//验证签名
		ok, err := wechat.VerifySign(constant.ApiKey, "HMAC-SHA256", request)
		if err != nil {
			fmt.Println("err:", err)
		}
		if ok != true {
			ctx.Text(replyMsg("FAIL", "签名失败"))
			return
			//panic(PaymentException.VerifySignError())
		}
		logUtils.Println("微信验签是否通过:", ok)
	}

	//println(bm)
	if err := request.ReturnCode; err != "SUCCESS" {
		logUtils.Println("微信通知接口接收通信失败:", request.ReturnMsg)
		ctx.Text(replyMsg("FAIL", ""))
		//panic(PaymentException.PaymentException())

	} else {
		var order db.OrderInfo
		if err := db.Driver.Where("number = ?", request.OutTradeNo).First(&order).Error;err==nil{
			logUtils.Println(orderException.OrderNotExsit())
			return
		}

		var payment db.WxPayOrder
		if err := db.Driver.Where("out_trade_no = ?", request.OutTradeNo).First(&payment).Error;err!=nil{
			payment.TradeStateDesc = request.ReturnMsg
			logUtils.Println(transactionException.PaymentException())
			return
		}

		if state := payment.TradeState; state == "SUCCESS" {
			logUtils.Println("回调已处理")
			return
		}
		if totalFee := request.TotalFee; util.String2Int(totalFee) != payment.TotalFee && util.String2Int(totalFee) != order.Total {
			logUtils.Println(transactionException.AmoutnIsNotEqual())
			ctx.Text(replyMsg("FAIL", "金额不等"))
			return
		}

		tx := db.Driver.Begin()
		order.Status = orderEnum.Done
		if err := tx.Debug().Save(&order).Error; err != nil {
			tx.Rollback()
			logUtils.Println(orderException.CancelRefuse())
			return
		}

		payment.Status = PaymentEnums.SUCCESS
		payment.TransactionId = request.TransactionId
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

	ctx.Text(replyMsg("SUCCESS", "OK"))

}

type NotifyResponse struct {
	ReturnCode string `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
}
func replyMsg(code, msg string) string {
	rsp := new(NotifyResponse)
	rsp.ReturnCode = code
	rsp.ReturnMsg = msg

	message, _ := json.Marshal(rsp)
	return string(message)
}

