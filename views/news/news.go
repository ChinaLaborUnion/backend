package news

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	homeworkException "grpc-demo/exceptions/homework"
	newsException "grpc-demo/exceptions/news"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

//curd 均需考虑是否存在级联操作,有就事务
//Create -- post  ctx就是类似于前端发送过来的请求体   auth 登陆态通过登陆态可以知道是谁登陆了这个app，学生，老师，管理员    原子操作：设计多表的curd需要用事务保证原子性
func CreateNews(ctx iris.Context,auth authbase.AuthAuthorization)  {
	//auth.CheckLogin()
	auth.CheckAdmin()
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	//判断标签是否存在，存在才给创建资讯
	label := params.Int("news_label","资讯标签")
	newsLabelIsExistById(label)

	title := params.Str("title","标题")
	content := params.Str("content","内容")
	isPublish := params.Bool("is_publish","是否发布")
	//更新时间,创建时间 不写，因为会自动更新
	news := db.NewsInfo{
		AccountId: auth.AccountModel().Id,
		Title: title,
		Content: content,
		IsPublish: isPublish,
		NewsLabelId: label,
	}
	if params.Has("introduction"){
		news.Introduction = params.Str("introduction","简介")
	}
	if params.Has("pictures"){
		picture := params.List("pictures","封面图片")
		var p string
		//目前理解是将一个东西x 序列化成为了 byte[],error，下面的dataPicture,_ 对应的是byte[],error;然后byte[] 需要转化为string,然后存进DB
		if dataPicture,err := json.Marshal(picture);err != nil{
			panic(homeworkException.PicturesMarshalFail())
		}else{
			p = string(dataPicture)
		}
		news.Pictures = p
		fmt.Println(news.Pictures)
	}
	db.Driver.Create(&news)

	ctx.JSON(iris.Map{
		"id":news.Id,
	})
}

// Put -- put
func PutNews(ctx iris.Context,auth authbase.AuthAuthorization,nid int){
	auth.CheckAdmin()
	//从缓存or数据库拿到的一条记录
	//根据传过来的newsId判断资讯是否存在
	var news db.NewsInfo
	if err:= db.Driver.GetOne("news_info",nid,&news);err != nil{
		panic(newsException.NewsNotExist())
	}
	////鉴权  创建者和管理员才能正常修改
	//if auth.AccountModel().Id != news.AccountId && !auth.IsAdmin(){
	//	//当前登陆者不能修改 别人的创建的班级
	//	panic(newsException.IllegalModify())
	//}

	//前端发来的 请求体
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(&news)
	//修改对应的数据
	//如果修改标签字段的话，需要根据表签名先判断是否存在
	if params.Has("news_label"){
		label := params.Int("news_label","资讯标签")
		newsLabelIsExistById(label)
		//标签存在的
		news.NewsLabelId = label
	}
	if params.Has("title"){
		news.Title = params.Str("title","标题")
	}
	if params.Has("introduction"){
		news.Introduction = params.Str("introduction","简介")
	}
	if params.Has("content"){
		news.Content = params.Str("content","内容")
	}
	if params.Has("is_publish"){
		news.IsPublish = params.Bool("is_publish","是否发布")
	}
	if params.Has("pictures"){
		picture := params.List("pictures","封面图片")
		var p string
		if dataPicture,err := json.Marshal(picture);err != nil{
			panic(homeworkException.PicturesMarshalFail())
		}else{
			p = string(dataPicture)
		}
		news.Pictures = p
	}
	//保存回去db
	db.Driver.Save(&news)
	//response 告诉前端 ok
	ctx.JSON(iris.Map{
		"id": news.Id,
	})
}

//Delete     ---  Delete     删除  cid int
func DeleteNews(ctx iris.Context,auth authbase.AuthAuthorization,nid int){
	auth.CheckAdmin()
	var news db.NewsInfo
	//从db or cache get one data
	// 拿数据
	//todo
	if err:= db.Driver.GetOne("news_info",nid,&news);err == nil{
		db.Driver.Delete(news)
	}
	////鉴权  创建者和管理员才可以delete
	//if auth.AccountModel().Id != news.AccountId && !auth.IsAdmin(){
	//	panic(newsException.IllegalDelete())
	//}else{
	//	//db delete
	//	db.Driver.Delete(news)
	//}
	//response
	ctx.JSON(iris.Map{
		"id":nid,
	})
}

//获取很多条记录  List --- get    需求：资讯所有人都可以看见，但是需要进行标签过滤（就像tb买东西时候，点击电器标签，下面的商品都是电器分类下的商品）
func ListNews(ctx iris.Context,auth authbase.AuthAuthorization){
	//自定义的list，很少个但是很关键的字段
	//这里定义的struct是table的字段的集合的子集  感觉像Model
	var lists []struct {
		Id         int   `json:"id"`
		CreateTime int64 `json:"create_time"`
	}
	//多少条记录
	var count int
	//specify the table on which you would like to run db operations
	table := db.Driver.Table("news_info")
	//一页多少条记录
	limit := ctx.URLParamIntDefault("limit", 10)
	//分页
	page := ctx.URLParamIntDefault("page", 1)
	//资讯所有人都可以看  所以不需要进行where
	//if !auth.IsAdmin() {
	//	table = table.Where("account_id = ?", auth.AccountModel().Id)
	//}
	//todo 根据标签过滤       done                  解析路径参数news_label_id 找不到 default值是 0
	if nlid := ctx.URLParamIntDefault("news_label_id",0); nlid != 0 {
		table = table.Where("news_label_id = ?",nlid)
	}
	//分页操作	//       ctx.URLParamIntDefault("x",0)
	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, create_time").Find(&lists)
	//向前端返回 lists
	ctx.JSON(iris.Map{
		"news": lists,
		"total": count,         //总共多少条
		"limit": limit,          //一页多少条
		"page":  page,        //当前页
	})
}

//todo  add some words  还要反序列化  done
var newsField = []string{
	"Id","Title","Introduction","Content","NewsLabelId","IsPublish","AccountId","CreateTime","UpdateTime",
}

//Mget --- post    根据前端给来的 id 数组，进行获取更详细的goods信息
func MgetNews(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.IsLogin()
	//因为在list接口的时候就已经按照身份进行get ids了，所以这里只要判断一下login就行
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")
	data := make([]interface{}, 0, len(ids))
	news := db.Driver.GetMany("news_info",ids,db.NewsInfo{})
	for _,newa := range news{
		func(data *[]interface{}){
			*data = append(*data,getData(newa.(db.NewsInfo)))
			defer func() {
				recover()
			}()
		}(&data)
	}
	//返回data
	ctx.JSON(data)
}

//反序列化
func getData(news db.NewsInfo)map[string]interface{}{
	v := paramsUtils.ModelToDict(news,newsField)
	var pictures []string
	if err := json.Unmarshal([]byte(news.Pictures),&pictures);err != nil{
		//fmt.Println("sbsbsbsb")
		panic(newsException.PicturesUnmarshalFail())
	}
	//json.Unmarshal([]byte(news.Pictures),&pictures)
	for _,i := range pictures{
		fmt.Println(i)
	}
	//因为是ModelToDict（Dictation所以就是picture）
	v["pictures"] = pictures
	return v
}

