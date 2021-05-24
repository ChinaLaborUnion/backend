package db

type WxPayOrder struct {
	Id int `gorm:"primary_key" json:"id"`

	Description string `gorm:"type:varchar(256)" json:"description"`

	Detail string `json:"detail"`

	AccountId int `json:"account_id" gorm:"not null;index"`

	OrderId int `json:"order_id" gorm:"not null;index"`

	//OpenId string `json:"open_id"`

	PrepayId string `json:"prepay_id"`

	TransactionId string `json:"transaction_id"`

	OutTradeNo string `json:"out_trade_no" gorm:"not null;index"`

	Status int `json:"status" gorm:"default:1"`

	TotalFee int `json:"total_fee"`

	DeviceInfo string `json:"device_info"`

	GoodsTag string `json:"goods_tag"`

	TradeType string `json:"trade_type" gorm:"default: 'MiniApp'"`

	TradeState string `json:"trade_state"`

	TradeStateDesc string `json:"trade_state_desc"`

	BankType string `json:"bank_type"`

	NonceStr string `json:"nonce_str"`

	PaySign string `gorm:"type:varchar(1024)" json:"pay_sign"`

	SpbilCreateIp string `json:"spbil_create_ip"`

	TimeStamp int64 `json:"time_stamp"`

	TimeStart int64 `json:"time_start"`

	TimeExpire int64 `json:"time_expire"`

	TimeEnd int64 `json:"time_end"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}

type AliPayOrder struct {
	Id int `gorm:"primary_key" json:"id"`
	//用户id
	AccountId int `json:"account_id"`
	//订单id
	OrderId int `json:"order_id"`
	//商品标题
	Subject string `json:"subject"`
	//商户订单号
	OutTradeNo string `json:"out_trade_no"`
	//支付宝交易流水号
	TradeNo string `json:"trade_no"`
	//收款支付宝账号对应的支付宝唯一用户号
	SellerId string `json:"seller_id"`
	//订单总金额
	TotalAmount string `json:"total_amount"`
	//商户请求参数签名串
	Sign string `json:"sign"`
	//网关返回码
	Code string `json:"code"`
	//返回码描述
	Msg string `json:"msg"`
	//支付宝交易状态
	TradeStatus string `json:"trade_status"`
	//后台状态
	Status int `json:"status" gorm:"default:1"`
	// 创建时间
	CreateTime int64 `json:"create_time"`

	TimeStamp int64 `json:"time_stamp"`
	// 更新时间
	UpdateTime int64 `json:"update_time"`

	TimeStart int64 `json:"time_start"`

	TimeExpire int64 `json:"time_expire"`

	TimeEnd int64 `json:"time_end"`
}