package courseException

import "grpc-demo/models"

func NotExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 6100,
		Message: "课程不存在",
	}
}

func PptNotExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 6200,
		Message: "PPT不存在",
	}
}

func VideoNotExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 6300,
		Message: "视频不存在",
	}
}

func GoodsNotExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 6300,
		Message: "商品不存在",
	}
}

func DoError() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 6300,
		Message: "执行失败",
	}
}

func PictureExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 6300,
		Message: "课程轮播图已存在",
	}
}

func PictureNotExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 6300,
		Message: "课程轮播图不存在",
	}
}
