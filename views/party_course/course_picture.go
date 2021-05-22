package party_course

import (
	"encoding/json"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	accountException "grpc-demo/exceptions/account"
	courseException "grpc-demo/exceptions/course"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreateCoursePicture(ctx iris.Context,auth authbase.AuthAuthorization){
	if !auth.IsAdmin(){
		panic(accountException.NoPermission())
	}

	// 是否存在记录
	table := db.Driver.Table("course_picture")
	var count int
	table.Count(&count)
	if count > 0{
		panic(courseException.PictureExist())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	picture := params.List("course_list","课程列表")

	//判断list是否有不存在的课程id
	var courses []db.PartyCourse
	db.Driver.Where("id in (?)",picture).Find(&courses)
	if len(courses) != len(picture){
		panic(courseException.NotExist())
	}

	var coursePicture db.CoursePicture

	data,_ := json.Marshal(picture)
	coursePicture.CourseList = string(data)

	db.Driver.Create(&coursePicture)
	ctx.JSON(iris.Map{
		"id":coursePicture.Id,
	})
}

func PutCoursePicture(ctx iris.Context,auth authbase.AuthAuthorization){
	if !auth.IsAdmin(){
		panic(accountException.NoPermission())
	}

	var coursePicture db.CoursePicture

	if err := db.Driver.Table("course_picture").First(&coursePicture).Error;err != nil{
		panic(courseException.PictureNotExist())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	if params.Has("course_list"){
		picture := params.List("course_list","课程列表")

		//判断list是否有不存在的课程id
		var courses []db.PartyCourse
		db.Driver.Where("id in (?)",picture).Find(&courses)
		if len(courses) != len(picture){
			panic(courseException.NotExist())
		}

		data,_ := json.Marshal(picture)
		coursePicture.CourseList = string(data)

		db.Driver.Save(&coursePicture)
	}

	ctx.JSON(iris.Map{
		"id":coursePicture.Id,
	})
}


func GetCoursePicture(ctx iris.Context){
	var coursePicture db.CoursePicture
	//
	//if err := db.Driver.GetOne("course_picture",pid,&coursePicture);err != nil{
	//	panic(courseException.PictureNotExist())
	//}

	if err := db.Driver.Table("course_picture").First(&coursePicture).Error;err != nil{
		panic(courseException.PictureNotExist())
	}

	var courseList []interface{}
	if err := json.Unmarshal([]byte(coursePicture.CourseList),&courseList);err != nil{
		panic(courseException.DoError())
	}

	data := make([]interface{}, 0, len(courseList))
	partyCourses := db.Driver.GetMany("party_course", courseList, db.PartyCourse{})
	for _,partyCourse  := range partyCourses {
		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(partyCourse,[]string{"Id","CourseCover"}))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}