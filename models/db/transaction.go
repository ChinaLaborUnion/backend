package db

//交易信息表
type TransactionInfo struct{
	//支付信息ID
	Id int `gorm:"primary_key" json:"id"`

	//订单ID
	OrderId int `json:"order_id" gorm:"index"`

	//用户ID
	AccountId int `json:"account_id"`

	//支付平台： 1-微信  2-支付宝
	Platform int `json:"platform"`

	//交易状态:  1-未付款 2-已付款 4-申请退款 8-已退款
	TransactionStatus int  `json:"transaction_status"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`

}

