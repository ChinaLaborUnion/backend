package order

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
)

func CreateOrder(ctx iris.Context,auth authbase.AuthAuthorization){
	//auth.CheckLogin()
	//params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	//goodsID := params.Int("goods_id","商品id")
	//var goods db.GoodsInfo
	//if err := db.Driver.GetOne("goods_info",goodsID,&goods);err != nil{
	//	panic(goodsException.GoodsNotExsit())
	//}
	//total := params.Int("total","总数")
	//
	////order.OrderNum = strconv.FormatInt(order.CreateTime, 10) + "-" + hash.GetRandomString(8)
}
