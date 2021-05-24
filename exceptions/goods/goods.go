package goodsException

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

func GoodsNotExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5400,
		Message: "商品不存在",
	}
}