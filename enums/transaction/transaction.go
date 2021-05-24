package transactionEnums

import enumsbase "grpc-demo/enums"


const (
	SUBMIT     = 1 //已提交
	PAID  	   = 2 //已付款
	REFUND     = 4 //提交退款
	REFUNDED   = 8 //已退款
	OVER       = 16 //订单结束
)

func NewOrderStatusEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 4, 8, 16},
	}
}
