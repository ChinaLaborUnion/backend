package wx

import (
	transactionEnums "grpc-demo/enums/transaction"
	"grpc-demo/views/transaction/payutil"

	"github.com/kataras/iris"
	"github.com/thinkeridea/go-extend/exnet"

	WxConstant "grpc-demo/constants/payment/wxpayment"
	PaymentEnums "grpc-demo/enums/payment"
	"grpc-demo/exceptions/transaction"
	"grpc-demo/models/db"
	"grpc-demo/utils/hash"
	logUtils "grpc-demo/utils/log"
	"strconv"
	"time"
)

var (
	appId  = WxConstant.AppId
	mchId  = WxConstant.MchId
)

func PaymentForWx(ctx iris.Context, oid, aid, total int, otn, body string) {

	notifyUrl := WxConstant.UnifiedNotifyUrl
	tradeType := "APP"
	scIp := exnet.ClientIP(ctx.Request())

	amount := make(map[string]interface{})
	amount["total"] = total
	amount["currency"] = "CNY"

	sceneInfo := make(map[string]string)
	sceneInfo["payer_client_ip"] = scIp

	paramMap := make(map[string]interface{})
	paramMap["appid"] = appId
	paramMap["mchid"] = mchId
	paramMap["description"] = body
	paramMap["out_trade_no"] = otn
	paramMap["notify_url"] = notifyUrl
	paramMap["amount"] = amount
	paramMap["scene_info"] = sceneInfo

	logUtils.Println(paramMap)

	nonceStr := hash.GetRandomString(32)
	Time := time.Now().Local()
	timeStamp := strconv.FormatInt(Time.Unix(), 10)
	timestamp2int := hash.String2Int64(timeStamp)

	//超时时间
	hr2, _ := time.ParseDuration("2h")
	timeAdd2 := Time.Add(hr2)
	timeStamp2 := strconv.FormatInt(timeAdd2.Unix(), 10)

	wxRsp, err := payutil.UnifiedOrder(paramMap)
	if err != nil {
		logUtils.Println("统一下单错误：", err)
		return
	}
	logUtils.Println(wxRsp["prepay_id"])

	appSign, err := payutil.AppSign(wxRsp, nonceStr, timeStamp)
	if err != nil {
		logUtils.Println("签名失败:", err)
	}
	logUtils.Println(appSign)


	transaction := db.TransactionInfo{
		OrderId:           oid,
		AccountId:         aid,
		Platform:          1,
		TransactionStatus: transactionEnums.SUBMIT,
	}
	tx := db.Driver.Begin()
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		panic(transactionException.TransactionCreateFail())
	}
	wxOrder := db.WxPayOrder{
		Status:      PaymentEnums.PENDING,
		Description: body,
		AccountId:   aid,
		OrderId:     oid,
		OutTradeNo:  otn,
		TotalFee:    total,
		TradeType:   tradeType,
		//PrepayId:    wxRsp["prepay_id"],
		//NonceStr:    nonceStr,
		//PaySign:       appSign,
		TimeStamp:     timestamp2int,
		TimeStart:     hash.String2Int64(time.Unix(timestamp2int, 0).Format("20060102150405")),
		TimeExpire:    hash.String2Int64(time.Unix(hash.String2Int64(timeStamp2), 0).Format("20060102150405")),
		SpbilCreateIp: scIp,
	}


	if err := tx.Create(&wxOrder).Error; err != nil {
		tx.Rollback()
		panic(transactionException.PaymentCreateFail())
	}

	tx.Commit()

	ctx.JSON(iris.Map{
		"request": iris.Map{
			"appid":     appId,
			"partnerid": mchId,
			"prepayid":  wxRsp["prepay_id"],
			"package":   "Sign=WXPay",
			"nonceStr":  nonceStr,
			"timestamp": timeStamp,
			"sign":      appSign,
		},
	})
}