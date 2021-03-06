package newsException

import "grpc-demo/models"

func PicturesMarshalFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "图片列表序列化失败",
	}
}

func PicturesUnmarshalFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "图片列表反序列化失败",
	}
}

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

