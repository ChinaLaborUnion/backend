package sign_up

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	signUpEnum "grpc-demo/enums/sign_up"
	accountException "grpc-demo/exceptions/account"
	signupException "grpc-demo/exceptions/signup"
	"grpc-demo/models/db"

	paramsUtils "grpc-demo/utils/params"
)

func CreatSignUp(ctx iris.Context,auth authbase.AuthAuthorization)  {
	auth.CheckLogin()
	userid := auth.AccountModel().Id
	//TODO 选课码-班
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	//classid := params.Int("class_id","班级id")
	cid := params.Int("course_id","course_id")
	//post方法写选课码
	code := params.Str("code","选课码")

	var s db.SignUp
	if err := db.Driver.Where("user_id = ? and course_id = ?",auth.AccountModel().Id,cid).First(&s).Error;err == nil{
		panic(signupException.SignedUp())
	}
	//var v db.Class
	var c db.PartyClass
	//通过选课码查找class1.id
	if err := db.Driver.Where("code = ? and party_course_id = ?",code,cid).First(&c).Error;err != nil{
		panic(signupException.SignupClassIdNotfound())
	}
	//classid = v.Id
	//if err := db.Driver.GetOne("class1",classcode,&c);err != nil{
	//		panic(signupException.SignupClassIdNotfound())
	//	}
	var signup db.SignUp
	signup = db.SignUp{
		//CourseId: v.CourseId,
		UserId: userid,
		ClassId: c.Id,
		CourseId: cid,
		Status: signUpEnum.NoDone,
	}
	db.Driver.Create(&signup)
	ctx.JSON(iris.Map{
		"id":signup.Id,
	})
}

//todo //error
func DeleteSignUp(ctx iris.Context,auth authbase.AuthAuthorization,sid int){
	//todo 权限——管理员/创建人可删除
	auth.CheckLogin()
	tid := auth.AccountModel().Id
	//判断管理员登录
	var SignUpUser db.SignUp
	if err := db.Driver.GetOne("sign_up",sid,&SignUpUser);err == nil{
		if !auth.IsAdmin() && tid != SignUpUser.UserId{
			panic(signupException.SignupUserNotHaveAuthority())
		}
		db.Driver.Delete(SignUpUser)
	}

	ctx.JSON(iris.Map{
		"id":sid,
	})
}

func PutSignUp(ctx iris.Context, cid int,auth authbase.AuthAuthorization) {
	//修改表 修改cid传入的人的状态
	auth.CheckLogin()
	//判断登录
	tid := auth.AccountModel().Id
	var signup db.SignUp
	if err := db.Driver.GetOne("sign_up", cid, &signup); err != nil {
		panic(signupException.SignupUsernotfound())
	}
	//tid对比account_id 为空——无权限 或 非管理员——无权限
	if !auth.IsAdmin() && tid != signup.UserId{
		panic(accountException.NoPermission())
	}
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(signup)
	//解释：params.Diff() 就是自己写的方法，如果前端传过来的请求体有这个字段，就修改，如果没有就从原来的这条记录拿。所以就不用if params.Has()
	//修改对应的数据
	if params.Has("status"){
		signup.Status = int16(params.Int("status","状态"))
	}
	statusEnums := signUpEnum.NewStatusEnums()
	if !statusEnums.Has(int(signup.Status)){
		panic(signupException.StatusIsNotAllow())
	}

	db.Driver.Save(&signup)
	ctx.JSON(iris.Map{
		"id":signup.Id,
	})

}

//todo
func SignUpListByCid(ctx iris.Context,auth authbase.AuthAuthorization,uid int)  {
	auth.CheckLogin()
	tid := auth.AccountModel().Id
	var Lists []struct{
		Id int `json:"id"`
		UserId int `json:"user_id"`
		ClassId int `json:"class_id"`
		CourseId int `json:"course_id"`
		Status int16 `json:"status"`
	}
	var count int
	var sg db.SignUp
	table := db.Driver.Table("sign_up")
	table = table.Where("class_id = ?",uid).First(&sg)
	if !auth.IsAdmin() && tid != sg.UserId{
		panic(signupException.SignupUserNotHaveAuthority())
	}

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)
	table.Count(&count).Offset((page - 1) * limit).Limit(limit).
		Select("id,user_id,class_id,course_id,status").Find(&Lists)
	ctx.JSON(iris.Map{
		"Lists" : Lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

func SignUpListByAid(ctx iris.Context,auth authbase.AuthAuthorization,aid int){
	auth.CheckLogin()
	tid := auth.AccountModel().Id
	var Lists []struct{
		Id int `json:"id"`
		UserId int `json:"user_id"`
		ClassId int `json:"class_id"`
		CourseId int `json:"course_id"`
		Status int16 `json:"status"`
	}
	var count int
	//TODO 判断登录
	table := db.Driver.Table("sign_up")
	table = table.Where("user_id = ?",aid)
	if !auth.IsAdmin() && tid != aid{
		panic(signupException.SignupUserNotHaveAuthority())
	}
	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)
	table.Count(&count).Offset((page - 1) * limit).Limit(limit).
		Select("id,user_id,class_id,course_id,status").Find(&Lists)

	ctx.JSON(iris.Map{
		"Lists" : Lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

func SignUpList(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckAdmin()
	var Lists []struct{
		Id int `json:"id"`
		UserId int `json:"user_id"`
		ClassId int `json:"class_id"`
		CourseId int `json:"course_id"`
		Status int16 `json:"status"`
	}
	var count int
	//TODO 判断登录
	table := db.Driver.Table("sign_up")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)
	table.Count(&count).Offset((page - 1) * limit).Limit(limit).
		Select("id,user_id,class_id,course_id,status").Find(&Lists)

	ctx.JSON(iris.Map{
		"Lists" : Lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

//func SignUpMegt(ctx iris.Context)  {
//	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
//	ids := params.List("ids", "id列表")
//	Registrant_data := make([]interface{}, 0, len(ids))
//	Registrants := db.Driver.GetMany("signup", ids, db.SignUp{})
//	for _, Registrant := range Registrants {
//		func(Registrant_data *[]interface{}) {
//			info := paramsUtils.ModelToDict(Registrant, []string{"Id","UserId",
//				"ClassId"})
//			*Registrant_data = append(*Registrant_data, info)
//			defer func() {
//				recover()
//			}()
//		}(&Registrant_data)
//	}
//
//	ctx.JSON(Registrant_data)
//}
