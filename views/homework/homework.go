package homework

import (
	"encoding/json"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	signUpEnum "grpc-demo/enums/sign_up"
	homeworkException "grpc-demo/exceptions/homework"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreateHomeWork(ctx iris.Context,auth authbase.AuthAuthorization,cid int){
	//todo 在url中加入班级id   done
	auth.CheckLogin()
	//本来需要判断班级存在不存在的，但是也没有必要判断。不在班级就是不在.
	//todo 如果此人不属于这个班级(学生没有加入班级，不允许上传作业) 不允许发布    done
	//根据班级id,和登陆者id，去找报名表，select * from signUp where class_id = cid and user_id = auth.AccountModel().Id
	var signUp db.SignUp
	if err:= db.Driver.Where("course_id = ? AND user_id = ?",cid,auth.AccountModel().Id).First(&signUp).Error;err != nil {
		//找不到
		panic(homeworkException.IllegalUpload())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	//todo 以下两个不需要在body中传     done    这两个东西都在 classId 就是cid ；courseId在那条记录里
	content := params.Str("content","内容")
	picture := params.List("picture","图片")


	var p string

	if dataPicture,err := json.Marshal(picture);err != nil{
		panic(homeworkException.PicturesMarshalFail())
	}else{
		p = string(dataPicture)
	}

	//todo 根据班级id 找到课程id 存在表中（适当冗余） done
	//好的，因为党课班级-学生是1-n，应该是根据班级Id去再次查找班级表的记录，然后去拿党课Id。适当冗余也不错，降到了数据库设计的设计原则中的第二范式。
	homework := db.Homework{
		//上传人
		UpperId : auth.AccountModel().Id,
		//班级id
		ClassId :signUp.ClassId,
		//课程id
		CourseId :signUp.CourseId,
		Content:content,
		//图片
		Picture : p,
		//视频

	}

	if params.Has("video"){
		video := params.List("video","视频")
		var v string
		if dataVideo,err := json.Marshal(video);err != nil{
			panic(homeworkException.VideosMarshalFail())
		}else{
			v = string(dataVideo)
		}
		homework.Video = v
	}

	tx := db.Driver.Begin()

	if err := tx.Create(&homework).Error;err != nil{
		tx.Rollback()
		panic(homeworkException.DoFail())
	}
	signUp.Status = signUpEnum.Done
	if err := tx.Save(&signUp).Error;err != nil{
		tx.Rollback()
		panic(homeworkException.DoFail())
	}
	tx.Commit()
	ctx.JSON(iris.Map{
		"id：":homework.Id,
	})
}

//改了
func PutHomeWork(ctx iris.Context,auth authbase.AuthAuthorization,hid int)  {
	auth.CheckLogin()

	var homework db.Homework
	//var classCreater db.Class1
	//accountId := auth.AccountModel().Id
	////查找创建者id
	if err := db.Driver.GetOne("homework",hid,&homework);err != nil{
		//这里的报错信息使用方法是:包名.类名
		panic(homeworkException.HomeworkNotExist())
	}
	//todo 卡权限 不是作业发布者不能修改
	//判断创建者id是否与登录者id吻合
	if homework.UpperId != auth.AccountModel().Id && !auth.IsAdmin(){
		panic(homeworkException.IllegalModify())
	}
	//前端发来的 请求体
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(&homework)
	//解释：params.Diff() 就是自己写的方法，如果前端传过来的请求体有这个字段，就修改，如果没有就从原来的这条记录拿。所以就不用if params.Has()
	//修改对应的数据
	//todo 以下三条不需要 也不能修改  done  属于不可修改字段

	if params.Has("picture"){
		var p string
		picture := params.List("picture","图片")
		if dataPicture,err := json.Marshal(picture);err != nil{
			panic(homeworkException.PicturesMarshalFail())
		}else{
			p = string(dataPicture)
		}
		homework.Picture = p
	}
	if params.Has("video"){
		var v string
		video := params.List("video","视频")
		if dataVideo,err := json.Marshal(video);err != nil{
			panic(homeworkException.VideosMarshalFail())
		}else{
			v = string(dataVideo)
		}
		homework.Video = v
	}
	if params.Has("content"){

		homework.Content = params.Str("content","内容")
	}
	db.Driver.Save(&homework)
	ctx.JSON(iris.Map{
		"id": homework.Id,
	})
}

//改了
func HomeWorkList(ctx iris.Context,auth authbase.AuthAuthorization)  {
	//todo 不需要在方法头中传cid   done
	//todo 直接在接口中允许在url中传cid 班级id，aid 账户id 进行过滤
	//todo 如果是 aid 账户id，如果是管理员 可以根据账户id过滤 否则本人只能看见自己的作业
	//todo 只传回ID和create_time              done
	auth.CheckLogin()
	var lists []struct {
		Id         int   `json:"id"`
		CreateTime int64 `json:"create_time"`
	}
	var count int
	//调用接口的角色分为用户[创建班级的老师和上传作业的学生]以及管理者       用户可以按照cid（班级id）过滤；管理员可以按照cid和aid（账号id）进行过滤
	table := db.Driver.Table("homework")
	//一页多少条记录
	limit := ctx.URLParamIntDefault("limit", 10)
	//分页
	page := ctx.URLParamIntDefault("page", 1)
	//班级id
	cid := ctx.URLParamIntDefault("cid",0)
	if cid != 0{
		var class db.PartyClass
		if err:= db.Driver.GetOne("party_class",cid,&class);err == nil{
			if class.AccountId == auth.AccountModel().Id || auth.IsAdmin(){
				table = table.Where("class_id = ?",cid)
			}else{
				table = table.Where("class_id = ? AND upper_id = ?",cid,auth.AccountModel().Id)
			}
		}else{
			table = table.Where("id = ?",0)
		}
	}

	//账号id
	aid := ctx.URLParamIntDefault("aid",0)
	if aid != 0{
		if !auth.IsAdmin(){
			table = table.Where("upper_id = ?",auth.AccountModel().Id)
		}else{
			table = table.Where("upper_id = ?",aid)
		}
	}
	if cid == 0 && aid == 0 && !auth.IsAdmin(){
		table = table.Where("upper_id = ?",auth.AccountModel().Id)
	}



	////管理员
	//if auth.IsAdmin() {
	//	if aid != 0{
	//		//select * from homework where aid = upper_id
	//		table = table.Where("upper_id = ?",aid)
	//	}
	//	if cid != 0{
	//		//select * from homework where cid = class_id
	//		table = table.Where("class_id = ?",cid)
	//	}
	//}else{
	//	//非管理员,如果是普通用户[包含创建班级的老师和学生]也传了aid，那就不管他，接下来就是cid班级.如果登陆者是创建cid的老师，那么在 加入班级signUp表中找到所有的学生
	//	var class db.PartyClass
	//	if err:= db.Driver.GetOne("party_class",cid,&class);err == nil {
	//		//拿到数据了，判断老师是不是这个班级的创建者
	//		if class.AccountId == auth.AccountModel().Id {
	//			//如果是老师调用这个接口，怎么办,从作业表里面拿，上传到该班级的作业都拿过来
	//			table = table.Where("class_id = ?",cid)
	//		}else {
	//			//学生查看自己某个班级下的作业
	//			table = table.Where("class_id = ? AND upper_id = ?",cid,auth.AccountModel().Id)
	//		}
	//	}else{
	//		panic(classException.ClassNotFount())
	//	}
	//}

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id,create_time").Find(&lists)
	ctx.JSON(iris.Map{
		"homeworks":  lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

//改了
func DeleteHomeWork(ctx iris.Context,auth authbase.AuthAuthorization,hid int)  {
	auth.CheckLogin()
	//todo 卡权限 不是作业发布者不能删除,   done：就是管理员或者作业发布者才可以delete
	var homework db.Homework
	if err:= db.Driver.GetOne("homework",hid,&homework);err == nil {
		//拿到该条记录
		if auth.AccountModel().Id == homework.UpperId || auth.IsAdmin(){
			db.Driver.Delete(homework)
		}else{
			panic(homeworkException.IllegalDelete())
		}
	}
	//todo 这里逻辑有问题 重写    done  重写如上
	ctx.JSON(iris.Map{
		"id":hid,
	})
}


func HomeWorkMegt(ctx iris.Context,auth authbase.AuthAuthorization)  {
	auth.CheckLogin()
	//因为在list接口的时候就已经按照身份进行get ids了，所以这里只要判断一下login就行
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")
	data := make([]interface{}, 0, len(ids))
	homeworks := db.Driver.GetMany("homework",ids,db.Homework{})
	for _,hw := range homeworks{
		var class db.PartyClass
		db.Driver.GetOne("party_class",hw.(db.Homework).ClassId,&class)
		if hw.(db.Homework).UpperId != auth.AccountModel().Id && !auth.IsAdmin() && class.AccountId != auth.AccountModel().Id{
			continue
		}
		func(data *[]interface{}){
			*data = append(*data,getData(hw.(db.Homework)))
			defer func() {
				recover()
			}()
		}(&data)
	}
	//返回data
	ctx.JSON(data)
}

var homeworkField = []string{
	"Id","UpperId","ClassId","CourseId","CreateTime","UpdateTime","Content",
}

//反序列化    Model
func getData(homework db.Homework)map[string]interface{}{
	v := paramsUtils.ModelToDict(homework,homeworkField)
	var pictures []string
	if homework.Video != ""{
		var videos []string
		if err := json.Unmarshal([]byte(homework.Video),&videos);err != nil{
			panic(homeworkException.VideosUnmarshalFail())
		}
		v["video"] = videos
	}

	if err := json.Unmarshal([]byte(homework.Picture),&pictures);err != nil{
		panic(homeworkException.PicturesUnmarshalFail())
	}

	//因为是ModelToDict（Dictation所以就是picture）
	v["picture"] = pictures

	return v
}


