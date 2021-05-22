package party_course

import (
	"encoding/json"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	accountException "grpc-demo/exceptions/account"
	"grpc-demo/exceptions/course"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreatePartyCourse(ctx iris.Context,auth authbase.AuthAuthorization){
	//auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	name := params.Str("party_course_name","课程名称")
	accountId := auth.AccountModel().Id
	courseBrief := params.Str("party_course_brief","课程简介")
	courseCover := params.Str("party_course_cover","课程封面")
	courseWork := params.Str("party_course_work","作业简介")

	partyCourse := db.PartyCourse{
		Name: name,
		AccountId: accountId,
		CourseBrief: courseBrief,
		CourseCover: courseCover,
		CourseWork: courseWork,
	}

	if params.Has("good_id"){
		goodId := params.Int("good_id","商品ID")

		//判断商品是否存在
		if err := db.Driver.GetOne("PartyCourse",goodId,&partyCourse);err == nil{
			panic(courseException.GoodsNotExist())
		}
		partyCourse.GoodsId = goodId
	}
	tx := db.Driver.Begin()

	if err := tx.Create(&partyCourse).Error;err != nil{
		tx.Rollback()
		panic(courseException.DoError())
	}

	var p db.PartyCourseData
	p.PartyCourseId = partyCourse.Id
	if params.Has("party_course_ppt"){
		ppt := params.List("party_course_ppt","课程ppt")
		data,_ := json.Marshal(ppt)
		p.CoursePpt = string(data)
	}

	if params.Has("party_course_video"){
		video := params.List("party_course_video","课程视频")
		data,_ := json.Marshal(video)
		p.CourseVideo = string(data)
	}


	if err := tx.Create(&p).Error;err != nil{
		tx.Rollback()
		panic(courseException.DoError())
	}


	tx.Commit()

	ctx.JSON(iris.Map{
		"id":partyCourse.Id,
	})
}

func PutPartyCourse(ctx iris.Context,auth authbase.AuthAuthorization,cid int){
	auth.CheckLogin()

	var partyCourse db.PartyCourse

	if err := db.Driver.GetOne("party_course",cid,&partyCourse);err != nil{
		panic(courseException.NotExist())
	}
	//判断权限
	if partyCourse.AccountId != auth.AccountModel().Id{
		panic(accountException.NoPermission())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(&partyCourse)
	partyCourse.Name = params.Str("party_course_name","课程名称")
	partyCourse.CourseBrief = params.Str("party_course_brief","课程简介")
	partyCourse.CourseCover = params.Str("party_course_cover","课程封面")
	partyCourse.CourseWork = params.Str("party_course_work","作业简介")

	if params.Has("good_id"){
		goodId := params.Int("good_id","商品ID")
		if err := db.Driver.GetOne("PartyCourse",goodId,&partyCourse);err != nil{
			panic(courseException.GoodsNotExist())
		}
		partyCourse.GoodsId = goodId
	}

	var partyCourseData db.PartyCourseData
	db.Driver.Where("party_course_id = ?",cid).First(&partyCourseData)
	if params.Has("party_course_ppt"){
		ppt := params.List("party_course_ppt","课程ppt")
		data,_ := json.Marshal(ppt)
		partyCourseData.CoursePpt = string(data)
	}

	if params.Has("party_course_video"){
		video := params.List("party_course_video","课程视频")
		data,_ := json.Marshal(video)
		partyCourseData.CourseVideo = string(data)
	}


	tx := db.Driver.Begin()
	if err := tx.Save(&partyCourseData).Error;err != nil{
		tx.Rollback()
		panic(courseException.DoError())
	}
	if err := tx.Save(&partyCourse).Error;err != nil{
		tx.Rollback()
		panic(courseException.DoError())
	}
	tx.Commit()

	ctx.JSON(iris.Map{
		"id":partyCourse.Id,
	})
}

func DeletePartyCourse(ctx iris.Context,cid int,auth authbase.AuthAuthorization){
	auth.CheckLogin()
	var partyCourse db.PartyCourse
	if err := db.Driver.GetOne("party_course",cid,&partyCourse);err == nil{
		if partyCourse.AccountId != auth.AccountModel().Id{
			panic(accountException.NoPermission())
		}
		db.Driver.Delete(partyCourse)
		db.Driver.Exec("delete from party_course_data where party_course_id = ?",cid)
	}

	ctx.JSON(iris.Map{
		"id":cid,
	})
}

func ListPartyCourse(ctx iris.Context){

	var lists []struct {
		Id int `json:"id"`
		CreateTime int64 `json:"create_time"`
	}

	var count int

	table := db.Driver.Table("party_course")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)


	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, create_time").Find(&lists)
	ctx.JSON(iris.Map{
		"party_course": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

func MgetPartyCourse(ctx iris.Context){
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	partyCourses := db.Driver.GetMany("party_course", ids, db.PartyCourse{})
	for _,partyCourse  := range partyCourses {
		func(data *[]interface{}) {
			*data = append(*data, GetData(partyCourse.(db.PartyCourse)))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}

var courseField = []string{
	"Id","Name","AccountId","CourseBrief","CourseCover","CourseWork","GoodsId","CreateTime","UpdateTime",
}

var resourceField = []string{
	"Id","PartyCourseId","CreateTime","UpdateTime",
}

func GetData(partyCourse db.PartyCourse)map[string]interface{}{
	data := paramsUtils.ModelToDict(partyCourse,courseField)
	var partyCourseDatas db.PartyCourseData
	db.Driver.Where("party_course_id = ?",partyCourse.Id).Find(&partyCourseDatas)

	p := paramsUtils.ModelToDict(partyCourseDatas,resourceField)
	var ppt []string
	json.Unmarshal([]byte(partyCourseDatas.CoursePpt),&ppt)
	p["ppt"] = ppt

	var video []string
	json.Unmarshal([]byte(partyCourseDatas.CourseVideo),&video)
	p["video"] = video

	data["data"] = p

	return data
}