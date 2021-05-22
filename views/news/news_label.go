package news

import (
	"fmt"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	newsException "grpc-demo/exceptions/news"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

//不存在的话就抛出异常
func newsLabelIsExistById(lid int){
	var newsLable db.NewsLabel
	if err:= db.Driver.GetOne("news_label",lid,&newsLable);err != nil{
		panic(newsException.NewsLableNotExist())
	}
}

func newsLabelIsExistByName(name string){
	var newsLable db.NewsLabel
	if err:= db.Driver.Where("name = ?",name).First(&newsLable).Error;err != nil{
		panic(newsException.NewsLableNotExist())
	}
}

//Create -- post  ctx就是类似于前端发送过来的请求体   auth 登陆态通过登陆态可以知道是谁登陆了这个app，学生，老师，管理员
func CreateNewsLabel(ctx iris.Context,auth authbase.AuthAuthorization)  {
	//todo 只有管理员可创建
	fmt.Println(auth.IsLogin())
	auth.CheckAdmin()
	//db.Driver.Debug().Exec("delete from news_label where id != 100")
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	//根据 前端发送过来的标签名，判断是否存在,如果不存在就可以创建
	name := params.Str("name","资讯标签")
    var nl db.NewsLabel
	if err:= db.Driver.Where("name = ?",name).First(&nl).Error;err != nil{
		//找不到这条记录，就去创建       更新时间,创建时间 不写，因为会自动更新
		newsLabel := db.NewsLabel{
			Name: name,
		}
		db.Driver.Create(&newsLabel)
		ctx.JSON(iris.Map{
			"id":newsLabel.Id,
		})
	}else{
		panic(newsException.NewsLableExist())
	}
}

// Put -- put
func PutNewsLabel(ctx iris.Context,auth authbase.AuthAuthorization,nlid int){
	auth.CheckAdmin()
	//根据传过来的news label id 判断是否存在这个标签
	var newsl db.NewsLabel
	var xx db.NewsLabel
	if err:= db.Driver.GetOne("news_label",nlid,&newsl);err != nil{
		panic(newsException.NewsLableNotExist())
	}
	//前端发来的 请求体
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	//判断传过来的标签名是否存在

	var name = params.Str("name","资讯标签名")
	if err:= db.Driver.Where("name = ?",name).First(&xx).Error;err != nil{
		//修改数值   就是找不到  然后再用getOne(id)  赋值这个newsl
		params.Diff(&newsl)
		newsl.Name = name
	}else{
		//err == nil  就是找到了  就是说传过来的标签名已经存在
		panic(newsException.NewsLableExist())
	}

	//保存回去db
	db.Driver.Save(&newsl)
	//response 告诉前端 ok
	ctx.JSON(iris.Map{
		"id": newsl.Id,
	})
}

//Delete     ---  Delete     需求：管理员可以删除。如果咨询表news1中的记录引用了即将被删除资讯标签的Id，则不允许操作
func DeleteNewsLabel(ctx iris.Context,auth authbase.AuthAuthorization,nlid int){
	auth.CheckAdmin()
	//todo 不需要 不存在不做操作  get it  done
	var newsLabel db.NewsLabel
	//从db or cache get one data
	var news db.NewsInfo
	//todo 如果有挂载 不允许操作

	if err := db.Driver.GetOne("news_label",nlid,&newsLabel);err == nil{
		if err:= db.Driver.Where("news_label_id = ?",nlid).First(&news).Error;err == nil{
			//err == nil 就是在news表里面找到了即将被删除的Id，所以news表中某一条记录引用了该标签Id    抛出非法删除异常
			panic(newsException.IllegalDelete())
		}
		db.Driver.Delete(newsLabel)
	}
	//response
	ctx.JSON(iris.Map{
		"id": nlid,
	})
}

//获取很多条记录  List --- get       需求：资讯标签谁都可以看
func ListNewsLabel(ctx iris.Context,auth authbase.AuthAuthorization){
	//自定义的list，很少个但是很关键的字段
	//这里定义的struct是table的字段的集合的子集  感觉像Model
	var lists []struct {
		Id         int   `json:"id"`
		Name string `json:"name"`
		CreateTime int64 `json:"create_time"`
	}
	//多少条记录
	var count int
	//specify the table on which you would like to run db operations
	table := db.Driver.Table("news_label")
	//一页多少条记录
	limit := ctx.URLParamIntDefault("limit", 10)
	//分页
	page := ctx.URLParamIntDefault("page", 1)

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id,name,create_time").Find(&lists)
	//向前端返回 lists
	ctx.JSON(iris.Map{
		"newsLabel": lists,
		"total": count,         //总共多少条
		"limit": limit,          //一页多少条
		"page":  page,        //当前页
	})
}