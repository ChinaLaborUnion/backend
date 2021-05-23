package order

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	goodsException "grpc-demo/exceptions/goods"
	"grpc-demo/models/db"
	"grpc-demo/utils/hash"
	paramsUtils "grpc-demo/utils/params"
	"strconv"
)

func CreateOrder(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	goodsID := params.Int("goods_id","商品id")
	var goods db.GoodsInfo
	if err := db.Driver.GetOne("goods_info",goodsID,&goods);err != nil{
		panic(goodsException.GoodsNotExsit())
	}
	total := params.Int("total","总数")
	totalPrice := total * goods.Price
	order := db.OrderInfo{
		Total: total,
		TotalPrice: totalPrice,
		AccountID: auth.AccountModel().Id,
		GoodsID: goodsID,
	}

	order.Number = strconv.FormatInt(order.CreateTime, 10) + "-" + hash.GetRandomString(8)
	db.Driver.Create(&order)
	ctx.JSON(iris.Map{
		"id":goods.ID,
	})
}

func PutOrders(ctx iris.Context,auth authbase.AuthAuthorization,oid int){
	auth.CheckAdmin()

	var order db.OrderInfo
	if err := db.Driver.GetOne("goods_info",oid,&order);err != nil{
		panic(goodsException.GoodsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(order)

	if params.Has("total"){
		total := params.Int("total","总数")
		var goods db.GoodsInfo
		if err := db.Driver.GetOne("goods_info",order.GoodsID,&goods);err != nil{
			panic(goodsException.GoodsNotExsit())
		}
		order.TotalPrice = total * goods.Price
	}

	db.Driver.Save(&order)
	ctx.JSON(iris.Map{
		"id":order.ID,
	})
}

func DeleteOrder(ctx iris.Context,auth authbase.AuthAuthorization,oid int){
	auth.CheckAdmin()

	var order db.OrderInfo
	if err := db.Driver.GetOne("goods_info",oid,&order);err == nil{
		db.Driver.Delete(&order)
	}

	ctx.JSON(iris.Map{
		"id":oid,
	})
}

func ListOrder(ctx iris.Context,auth authbase.AuthAuthorization){
	var lists []struct{
		ID         int   `json:"id"`
		CreateTime int64 `json:"create_time"`
	}
	var count int

	table := db.Driver.Table("order")

	limit := ctx.URLParamIntDefault("limit", 10)

	page := ctx.URLParamIntDefault("page", 1)

	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, create_time").Find(&lists)

	ctx.JSON(iris.Map{
		"orders": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

var goodsField = []string{"ID","Number","GoodsID","AccountID","Total","TotalPrice","CreateTime","UpdateTime"}

func MgetOrders(ctx iris.Context,auth authbase.AuthAuthorization){
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	orders := db.Driver.GetMany("order_info",ids,db.OrderInfo{})
	for range orders{
		func(data *[]interface{}){
			*data = append(*data,paramsUtils.ModelToDict(orders,goodsField))
			defer func() {
				recover()
			}()
		}(&data)
	}

	ctx.JSON(data)
}

