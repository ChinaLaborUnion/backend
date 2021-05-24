package ali

import (
	aliconstant "grpc-demo/constants/payment/alipayment"
	transactionEnums "grpc-demo/enums/transaction"
	transactionException "grpc-demo/exceptions/transaction"
	"grpc-demo/models/db"
	util "grpc-demo/utils/hash"
	logUtils "grpc-demo/utils/log"
	"strconv"
	"time"

	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/alipay"
	"github.com/kataras/iris"
)

var (
	client     *alipay.Client
	appId      = aliconstant.AppId
	privateKey = aliconstant.PrivateKey
)

func init() {

	time.FixedZone("CST", 3600*8)

	// 初始化支付宝客户端
	//    appId：应用ID
	//    privateKey：应用私钥，支持PKCS1和PKCS8
	//    isProd：是否是正式环境
	client = alipay.NewClient(appId, privateKey, true)
	// 设置国家，不设置默认就是 China
	client.
		SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetPrivateKeyType(alipay.PKCS1).
		SetNotifyUrl(aliconstant.NotifyUrl)

	err := client.SetCertSnByPath(aliconstant.AppCertPublicKey, aliconstant.AliPayRootCert, aliconstant.AppCertPublicKeyRSA2)
	if err != nil {
		logUtils.Println("SetCertSnByPath:", err)
		return
	}
	//os.Exit(m.Run())
}
func PaymentForAli(ctx iris.Context, oid, aid int, otn, subject, total string) {
	bm := make(gopay.BodyMap)
	//校验数据
	Time := time.Now()
	timeStamp := strconv.FormatInt(Time.Unix(), 10)
	timestamp2int := util.String2Int64(timeStamp)
	timestamp := util.GetTimeString2(Time)
	logUtils.Println(timestamp)

	hr2, _ := time.ParseDuration("2h")
	timeAdd2 := Time.Add(hr2)
	timeexpire := util.GetTimeString2(timeAdd2)
	timeStamp2 := strconv.FormatInt(timeAdd2.Unix(), 10)
	logUtils.Println(timeexpire)
	bm.Set("subject", subject)
	bm.Set("out_trade_no", otn)
	bm.Set("total_amount", total)
	bm.Set("timestamp", timestamp)
	bm.Set("time_expire", timeexpire)

	//调用APP支付接口2.0
	payParam, err := client.TradeAppPay(bm)
	if err != nil {
		logUtils.Println("err:", err)
		return
	}
	logUtils.Println("payParam:", payParam)

	var aliOrder = db.AliPayOrder{
		AccountId:  aid,
		OrderId:    oid,
		OutTradeNo: otn,
		Subject:    subject,
		TimeStamp:  timestamp2int,
		TimeStart:  util.String2Int64(time.Unix(timestamp2int, 0).Format("20060102150405")),
		TimeExpire: util.String2Int64(time.Unix(util.String2Int64(timeStamp2), 0).Format("20060102150405")),
	}

	transaction := db.TransactionInfo{
		OrderId:           oid,
		AccountId:         aid,
		Platform:          2,
		TransactionStatus: transactionEnums.SUBMIT,
	}
	tx := db.Driver.Begin()
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		panic(transactionException.TransactionCreateFail())
	}

	if err := tx.Create(&aliOrder).Error; err != nil {
		tx.Rollback()
		panic(transactionException.PaymentCreateFail())
	}
	tx.Commit()

	//传回[]byte类型数据
	ctx.Write([]byte(payParam))
}
