package newsException

import "grpc-demo/models"

func NewsLableNotExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "标签不存在",
		Data:    nil,
	}
}

func NewsLableExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "标签已存在",
		Data: nil,
	}
}

func NewsNotExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "资讯不存在",
		Data:    nil,
	}
}


func IllegalModify() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "非法修改资讯（资讯标签）",
	}
}

func IllegalDelete() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "非法删除资讯（资讯标签）",
	}
}

