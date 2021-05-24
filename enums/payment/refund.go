package payment

import enumsbase "grpc-demo/enums"

const (
	REFUNDSUCCESS		= 1 //退款成功
	REFUNDCLOSE 		= 2 //退款关闭。
	PROCESSING			= 4	//退款处理中
	CHANGE				= 8  /*	退款异常，退款到银行发现用户的卡作废或者冻结了，
								导致原路退款银行卡失败，可前往商户平台（pay.weixin.qq.com）-交易中心，
								手动处理此笔退款。$n为下标，从0开始编号。*/
)

func NewRefundEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 4, 8},
	}
}