package homework

import (
	"encoding/json"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	classException "grpc-demo/exceptions/class"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreatHomeWork(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()
	upperId := auth.AccountModel().Id
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	classId := params.Int("classId","班级id")
	courseId := params.Int("courseId","课程id")
	picture := params.List("picture","图片")
	video := params.List("video","视频")
	dataPicture,_ := json.Marshal(picture)
	dataVideo,_ := json.Marshal(video)
	homework := db.HomeWork{
		//上传人
		UpperId : upperId,
		//班级id
		ClassId :classId,
		//课程id
		CourseId :courseId,
		//图片
		Picture : string(dataPicture),
		//视频
		Video : string(dataVideo),
	}
	db.Driver.Create(&homework)
	ctx.JSON(iris.Map{
		"id：":homework.Id,
	})
}

func PutHomeWork(ctx iris.Context,auth authbase.AuthAuthorization,cid int)  {
	auth.CheckLogin()
	var homework db.HomeWork
	//var classCreater db.Class1
	//accountId := auth.AccountModel().Id
	////查找创建者id
	if err := db.Driver.GetOne("homework",cid,&homework);err != nil{
		//这里的报错信息使用方法是:包名.类名
		panic(classException.ClassNotFount())
	}
	////判断创建者id是否与登录者id吻合
	//if accountId != classCreater.AccountId && !auth.IsAdmin(){
	//	//当前登陆者不能修改 别人的创建的班级
	//	panic(classException.IllegalModify())
	//}
	//前端发来的 请求体

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(&homework)
	//解释：params.Diff() 就是自己写的方法，如果前端传过来的请求体有这个字段，就修改，如果没有就从原来的这条记录拿。所以就不用if params.Has()
	//修改对应的数据
	if params.Has("upper_id"){
		homework.UpperId = params.Int("upper_id","上传者")
	}
	if params.Has("class_id"){
		homework.ClassId = params.Int("class_id","班级号")
	}
	if params.Has("course_id"){
		homework.CourseId = params.Int("course_id","课程号")
	}
	if params.Has("picture"){
		picture := params.List("picture","图片")
		dataPicture,_ := json.Marshal(picture)
		homework.Picture = string(dataPicture)
	}
	if params.Has("video"){
		video := params.List("video","视频")
		dataVideo,_ := json.Marshal(video)
		homework.Picture = string(dataVideo)
	}
	db.Driver.Save(&homework)
	ctx.JSON(iris.Map{
		"id": homework.Id,
	})
}

func HomeWorkList(ctx iris.Context,auth authbase.AuthAuthorization,cid int)  {
	//todo cid课程号，去查找这门课的作业，所有人对作业可见
	auth.CheckLogin()
	var lists []struct {
		Id         int   `json:"id"`
		Picture        int   `json:"picture"`
		Video int64 `json:"video"`
	}
	var count int
	table := db.Driver.Table("home_work").Where("course_id = ?",cid)
	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)
	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id,picture,video").Find(&lists)
	ctx.JSON(iris.Map{
		"likes":  lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

func DeleteHomeWork(ctx iris.Context,auth authbase.AuthAuthorization,cid int)  {
	auth.CheckLogin()
	var homeWork db.HomeWork
	//todo 判断登录状态，用登录者id在class1表中查找账号id，如果非创建者账号，或非管理员，报错
	if err := db.Driver.Table("class1").Where("account_id = ?",auth.AccountModel().Id);err == nil || !auth.IsAdmin() {
		panic("无权限")
	}
	if err := db.Driver.Table("home_work").Where("id = ?",cid);err == nil{
		//成功拿到这条记录
		//todo 判断登陆者是不是创建者   done
		db.Driver.Delete(homeWork)
	}
	//response
	ctx.JSON(iris.Map{
		"id":cid,
	})
}

func HomeWorkMegt(ctx iris.Context,auth authbase.AuthAuthorization,cid int)  {
	auth.CheckLogin()
	uid :=auth.AccountModel().Id
	//type data struct {
	//	Ids []int `json:"ids"`
	//}
	var homework []db.HomeWork
	db.Driver.Where("account_id = ?", uid).Find(&homework)
	homeworkData := make([]interface{}, 0, len(homework))
	for _, hw := range homework {
		func(data *[]interface{}) {
			info := paramsUtils.ModelToDict(hw, []string{"Id", "UpperId", "ClassId","CourseId","Picture",
				"Video"})
			*data = append(*data, info)
			defer func() {
				recover()
			}()
		}(&homeworkData)
	}
	ctx.JSON(iris.Map{
		"data":homeworkData,
	})
}

