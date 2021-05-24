package transactionException

import "grpc-demo/models"

func TransactionCreateFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "交易创建失败",
	}
}
func TransactionGetFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "查无交易",
	}
}

func TransactionUpdateFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "交易更改失败",
	}
}

func PaymentException() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "付款异常",
	}
}
func PaymentCreateFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "付款订单创建失败",
	}
}
func PaymentSaveFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "付款订单创建失败",
	}
}

func PaymentCloseException() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "订单关闭异常",
	}
}
func PaymentQuerryException() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "订单查询异常",
	}
}

func AmoutnIsNotEqual() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "金额不对等",
	}
}

func VerifySignError() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "验证签名不通过",
	}
}

func UnifiedOrderFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "统一下单失败",
	}
}

func CanNotFind() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "无法找到支付单",
	}
}
func RefundCreateFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "退款创建失败",
	}
}

func QueryRefundFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "查询退款失败",
	}
}

func StatusIsNotAllow() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5100,
		Message: "状态不允许",
	}
}
