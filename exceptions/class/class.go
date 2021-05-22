package classException

import "grpc-demo/models"

func ClassNotFount() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "班级不存在",
	}
}

func IllegalModify() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "非法修改班级",
	}
}

func IllegalDelete() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "非法删除班级",
	}
}

