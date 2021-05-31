package wxv2

import (
	"encoding/json"
	"fmt"
	constant "grpc-demo/constants/payment/wxpayment"
	PaymentEnums "grpc-demo/enums/payment"
	transactionEnums "grpc-demo/enums/transaction"
	transactionException "grpc-demo/exceptions/transaction"
	"grpc-demo/models/db"
	"grpc-demo/utils/hash"
	logUtils "grpc-demo/utils/log"
	"strconv"
	"time"

	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/kataras/iris"
	"github.com/thinkeridea/go-extend/exnet"
)

var (
	client *wechat.Client
	appId  = constant.AppId
	mchId  = constant.MchId
	apiKey = constant.ApiKey
)

func init() {
	// 初始化微信客户端
	//    appId：应用ID
	//    mchId：商户ID
	//    apiKey：API秘钥值
	//    isProd：是否是正式环境
	client = wechat.NewClient(appId, mchId, apiKey, true)
	// 设置国家，不设置默认就是 China
	client.SetCountry(wechat.China)
	client.AddCertFilePath(constant.ApiClientCertPem, constant.ApiClientKey, constant.ApiClientCertP12)
}
func PaymentForWxV2(ctx iris.Context, oid, aid, total int, otn, body string) {

	notifyUrl := constant.UnifiedNotifyUrl
	tradeType := wechat.TradeType_App
	scIp := exnet.ClientIP(ctx.Request())

	bm := make(gopay.BodyMap)


	//校验数据
	bm.Set("sign_type", wechat.SignType_HMAC_SHA256)
	bm.Set("nonce_str", hash.GetRandomString(16))
	bm.Set("body", body)
	bm.Set("out_trade_no", otn)
	bm.Set("total_fee", total)
	bm.Set("spbill_create_ip", scIp) //exnet.ClientIP(ctx.Request())
	bm.Set("notify_url", notifyUrl)
	bm.Set("trade_type", tradeType)

	// 正式
	_b, _ := json.Marshal(bm)
	fmt.Println(string(_b))
	sign1 := wechat.GetParamSign(appId, mchId, apiKey, bm)

	fmt.Println("sign")
	fmt.Println(sign1)
	// 沙箱
	/*
		sign, err := wxpay.GetSanBoxParamSign("wx1328c016e69fdf9f", "1601892301", "9h8Aiuqh81g87f2bGwg567rF387GFR6f", bm)
		fmt.Println("sign")
		fmt.Println(err)
		fmt.Println(sign)
	*/

	bm.Set("sign", sign1)

	//请求支付下单，成功后得到结果
	//if wxRsp, err := client.UnifiedOrder(bm);err !=nil{
	//	panic(exception.PaymentException())
	//}
	wxRsp, err := client.UnifiedOrder(bm)
	fmt.Println(wxRsp)
	if err != nil {
		fmt.Println(err)
		logUtils.Println("统一下单错误:", err)
		panic(err)
	}
	Time := time.Now()
	timeStamp := strconv.FormatInt(Time.Unix(), 10)
	timestamp2int := hash.String2Int64(timeStamp)

	//超时时间
	hr2, _ := time.ParseDuration("2h")
	timeAdd2 := Time.Add(hr2)
	timeStamp2 := strconv.FormatInt(timeAdd2.Unix(), 10)


	nonceStr := hash.GetRandomString(16)
	paySign := wechat.GetAppPaySign(appId, mchId, nonceStr, wxRsp.PrepayId, wechat.SignType_HMAC_SHA256, timeStamp, apiKey)

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

	var wxOrder = db.WxPayOrder{
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

	if err := tx.Debug().Create(&wxOrder).Error; err != nil {
		tx.Rollback()
		panic(transactionException.PaymentCreateFail())
	}


	//tx.Debug().Save(&order)
	tx.Commit()

	ctx.JSON(iris.Map{
		"request": iris.Map{
			"appid":     appId,
			"partnerid": mchId,
			"prepayid":  wxRsp.PrepayId,
			"package":   "Sign=WXPay",
			"nonceStr":  nonceStr,
			"timestamp": timeStamp,
			"sign":      paySign,
		},
	})

}