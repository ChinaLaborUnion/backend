package payment

import enumsbase "grpc-demo/enums"

const (
	PENDING 		= 1     //待支付
	SUCCESS 		= 2		//支付成功
	FAIL 			= 4		//支付失败
	CANCEL 			= 8		//交易取消
	REFUNDING 		= 16	//退款中
	REFUND 			= 32	//退款成功
	CLOSED 			= 64	//交易超时关闭
)

func NewPaymentEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 4, 8, 16, 32, 64},
	}
}
