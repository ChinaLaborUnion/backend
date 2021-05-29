package accountException

import (
"grpc-demo/models"
)

func AuthIsNotLogin() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5300,
		Message: "尚未登录",
	}
}

func NoPermission() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5301,
		Message: "无权限执行此操作",
	}
}

func AccountNotFount() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5301,
		Message: "账户不存在",
	}
}

func EmailSendFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5301,
		Message: "验证码发送失败",
	}
}

func RedisFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5301,
		Message: "缓存操作失败",
	}
}

func ValidatedFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5301,
		Message: "验证码错误",
	}
}

func EmailValidatedFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5301,
		Message: "邮箱格式错误",
	}
}

func EmailRepeated() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5301,
		Message: "邮箱已存在",
	}
}

