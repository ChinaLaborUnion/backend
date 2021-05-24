package order

import (
	authbase "grpc-demo/core/auth"
	orderEnum "grpc-demo/enums/order"
	accountException "grpc-demo/exceptions/account"
	goodsException "grpc-demo/exceptions/goods"
	orderException "grpc-demo/exceptions/order"
	"grpc-demo/models/db"
	"grpc-demo/utils/hash"
	paramsUtils "grpc-demo/utils/params"
	"strconv"

	"github.com/kataras/iris"
)

func CreateOrder(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	goodsID := params.Int("goods_id","商品id")
	var goods db.GoodsInfo
	if err := db.Driver.GetOne("goods_info",goodsID,&goods);err != nil{
		panic(goodsException.GoodsNotExist())
	}
	total := params.Int("total","总数")
	if goods.Inventory < total{
		panic(orderException.InventoryNoEnough())
	}
	totalPrice := total * goods.Price
	order := db.OrderInfo{
		Status:orderEnum.WaitToPay,
		Total: total,
		TotalPrice: totalPrice,
		AccountID: auth.AccountModel().Id,
		GoodsID: goodsID,
	}
	tx := db.Driver.Begin()
	if err := tx.Create(&order).Error;err != nil{
		tx.Rollback()
		panic(orderException.CreateFail())
	}

	order.Number = strconv.FormatInt(order.CreateTime, 10) + "-" + hash.GetRandomString(8)
	if err := tx.Save(&order).Error;err != nil{
		tx.Rollback()
		panic(orderException.CreateFail())
	}
	goods.Inventory -= total
	if err := tx.Save(&goods).Error;err != nil{
		tx.Rollback()
		panic(orderException.CreateFail())
	}
	tx.Commit()
	ctx.JSON(iris.Map{
		"id":order.ID,
	})
}


func CancelOrder(ctx iris.Context,auth authbase.AuthAuthorization,oid int){
	auth.CheckLogin()

	var order db.OrderInfo
	if err := db.Driver.GetOne("order_info", oid, &order); err != nil {
		panic(orderException.OrderNotExsit())
	}

	if order.AccountID != auth.AccountModel().Id {
		panic(accountException.NoPermission())
	}


	if order.Status != orderEnum.WaitToPay{
		panic(orderException.CancelRefuse())
	}

	order.Status = orderEnum.Cancel
	tx := db.Driver.Begin()
	if err := tx.Save(&order).Error;err != nil{
		tx.Rollback()
		panic(orderException.CancelFail())
	}
	var goods db.GoodsInfo
	if err := db.Driver.GetOne("goods_info",order.GoodsID,&goods);err != nil{
		panic(goodsException.GoodsNotExist())
	}
	goods.Inventory += order.Total
	if err := tx.Save(&goods).Error;err != nil{
		tx.Rollback()
		panic(orderException.CancelFail())
	}
	tx.Commit()


	ctx.JSON(iris.Map{
		"id": order.ID,
	})
}

func ListOrder(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()
	var lists []struct{
		ID         int   `json:"id"`
		CreateTime int64 `json:"create_time"`
	}
	var count int

	table := db.Driver.Table("order_info")

	limit := ctx.URLParamIntDefault("limit", 10)

	page := ctx.URLParamIntDefault("page", 1)

	if !auth.IsAdmin() {
		table = table.Where("account_id = ?", auth.AccountModel().Id)
	}

	if author := ctx.URLParamIntDefault("author_id", 0); author != 0 && auth.IsAdmin() {
		table = table.Where("account_id = ?", author)
	}

	if status := ctx.URLParamIntDefault("status", 0); status != 0 {
			table = table.Where("status = ?", status)
	}

	if startTime := ctx.URLParamInt64Default("start_time", 0); startTime != 0 {
		endTime := ctx.URLParamInt64Default("end_time", 0)
		table = table.Where("create_time between ? and ?", startTime, endTime)
	}

	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, create_time").Find(&lists)

	ctx.JSON(iris.Map{
		"orders": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

var orderField = []string{"ID","Number","GoodsID","AccountID","Total","TotalPrice","Status","CreateTime","UpdateTime"}

func MgetOrders(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	orders := db.Driver.GetMany("order_info",ids,db.OrderInfo{})
	for _,o := range orders{
		if o.(db.OrderInfo).AccountID != auth.AccountModel().Id && !auth.IsAdmin() {
			continue
		}
		func(data *[]interface{}){
			*data = append(*data,paramsUtils.ModelToDict(o,orderField))
			defer func() {
				recover()
			}()
		}(&data)
	}

	ctx.JSON(data)
}

