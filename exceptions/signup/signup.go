package signupException

import "grpc-demo/models"

func SignupUsernotfound() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 7069,
		Message: "找不到该登记人",
	}
}

func SignupUserNotHaveAuthority() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 7075,
		Message: "没有权限",
	}
}

func SignupClassIdNotfound() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 7070,
		Message: "找不到该Id",
	}
}

func SignupNotfound() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 7071,
		Message: "找不到对象",
	}
}

func StatusIsNotAllow() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 7072,
		Message: "状态类型不允许",
	}
}

func SignedUp() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 7070,
		Message: "已经报名",
	}
}