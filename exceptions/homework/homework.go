package homeworkException

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

func VideosMarshalFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "视频列表序列化失败",
	}
}

func VideosUnmarshalFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "视频列表反序列化失败",
	}
}

func IllegalModify() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "非法修改",
	}
}

func IllegalDelete() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "非法删除",
	}
}

func IllegalUpload() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "尚未报名该班级，不可上传作业",
		Data:    nil,
	}
}

func HomeworkNotExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "作业不存在",
		Data:    nil,
	}
}
