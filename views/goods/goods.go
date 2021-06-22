package goods

import (
	"encoding/json"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	goodsException "grpc-demo/exceptions/goods"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreateGoods(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	cover := params.Str("cover","封面")
	name := params.Str("name","商品名")
	price := params.Int("price","价格")
	pictures := params.List("pictures","图片列表")
	inventory := params.Int("inventory","库存")

	var p string
	if data,err := json.Marshal(pictures);err != nil{
		panic(goodsException.PicturesMarshalFail())
	}else{
		p = string(data)
	}
	goods := db.GoodsInfo{
		Cover:      cover,
		Name:       name,
		Price:      price,
		Pictures:   p,
		Inventory:inventory,
	}
	if params.Has("brief"){
		goods.Brief = params.Str("brief","简介")
	}
	if params.Has("is_on"){
		goods.IsOn = params.Bool("is_on","是否上线")
	}

	db.Driver.Create(&goods)
	ctx.JSON(iris.Map{
		"id":goods.ID,
	})

}

func PutGoods(ctx iris.Context,auth authbase.AuthAuthorization,gid int){
	auth.CheckAdmin()

	var goods db.GoodsInfo
	if err := db.Driver.GetOne("goods_info",gid,&goods);err != nil{
		panic(goodsException.GoodsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(goods)
	goods.Name = params.Str("name","商品名")
	goods.Cover = params.Str("cover","封面")
	goods.Price = params.Int("price","价格")
	goods.Brief = params.Str("brief","简介")
	goods.IsOn = params.Bool("is_on","是否上线")
	goods.Inventory = params.Int("inventory","库存")

	if params.Has("pictures"){
		pictures := params.List("pictures","图片列表")
		var p string
		if data,err := json.Marshal(pictures);err != nil{
			panic(goodsException.PicturesMarshalFail())
		}else{
			p = string(data)
		}
		goods.Pictures = p
	}

	db.Driver.Save(&goods)
	ctx.JSON(iris.Map{
		"id":goods.ID,
	})
}

func DeleteGoods(ctx iris.Context,auth authbase.AuthAuthorization,gid int){
	auth.CheckAdmin()

	var goods db.GoodsInfo
	if err := db.Driver.GetOne("goods_info",gid,&goods);err == nil{
		goods.IsOn = false
		db.Driver.Save(&goods)
	}

	ctx.JSON(iris.Map{
		"id":gid,
	})
}

func ListGoods(ctx iris.Context,auth authbase.AuthAuthorization){
	var lists []struct{
		ID         int   `json:"id"`
		CreateTime int64 `json:"create_time"`
	}
	var count int

	table := db.Driver.Table("goods_info")

	limit := ctx.URLParamIntDefault("limit", 10)

	page := ctx.URLParamIntDefault("page", 1)

	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).
		Limit(limit).Select("id, create_time").Find(&lists)

	ctx.JSON(iris.Map{
		"goods": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

var goodsField = []string{"ID","Cover","Name","Brief","Price","People","IsOn","Inventory","CreateTime","UpdateTime"}

func MgetGoods(ctx iris.Context,auth authbase.AuthAuthorization){
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	goods := db.Driver.GetMany("goods_info",ids,db.GoodsInfo{})
	for _,v := range goods{
		func(data *[]interface{}){
			*data = append(*data,getData(v.(db.GoodsInfo)))
			defer func() {
				recover()
			}()
		}(&data)
	}

	ctx.JSON(data)
}

func getData(goods db.GoodsInfo)map[string]interface{}{
	v := paramsUtils.ModelToDict(goods,goodsField)
	var pictures []string
	if err := json.Unmarshal([]byte(goods.Pictures),&pictures);err != nil{
		panic(goodsException.PicturesUnmarshalFail())
	}
	v["pictures"] = pictures
	return v
}




