package party_class

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	classException "grpc-demo/exceptions/class"
	"grpc-demo/exceptions/course"
	"grpc-demo/models/db"
	"grpc-demo/utils/hash"
	paramsUtils "grpc-demo/utils/params"
)

//Create -- post  ctx就是类似于前端发送过来的请求体   auth 登陆态通过登陆态可以知道是谁登陆了这个app，学生，老师，管理员
func CreatePartyClass(ctx iris.Context,auth authbase.AuthAuthorization,pid int)  {
	auth.CheckLogin()
	//todo 党课是否存在  done
	var partyCourse db.PartyCourse
	if err := db.Driver.GetOne("party_course",pid,&partyCourse);err != nil{
		//这里的报错信息使用方法是:包名.类名
		panic(courseException.NotExist())
	}
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	name := params.Str("name","班级名称")

	class := db.PartyClass{
		Name:      name,
		PartyCourseId: partyCourse.Id,
	}
	if params.Has("introduce") {
		class.Introduce = params.Str("introduce","班级简介")
	}
	//这里使用utils包下的生成长度为6的随机序列,没有去check 班级码的唯一性
	class.Code = hash.GetRandomInt(6)
	if params.Has("place"){
		class.Place = params.Str("place","地点")
	}
	//这里前端传过来的就是int64的时间戳，但是我们写的params.Int()是转化为int的，所以我们还需要进行类型转换。
	if params.Has("start_time") {
		class.StartTime = int64(params.Int("start_time","开始时间"))
	}
	if params.Has("end_time") {
		class.EndTime = int64(params.Int("end_time","结束时间"))
	}
	class.AccountId = auth.AccountModel().Id
	class.TeacherName = params.Str("teacher_name","教师名字")
	//更新时间,创建时间 不写，因为会自动更新
	if params.Has("comment") {
		class.Comment = params.Str("comment","备注")
	}
	db.Driver.Create(&class)

	ctx.JSON(iris.Map{
		"id":class.Id,
	})
}


// Put -- put                   修改一条班级记录
func PutPartyClass(ctx iris.Context,auth authbase.AuthAuthorization,cid int){
	auth.CheckLogin()
	//从缓存or数据库拿到的一条活动记录
	//根据这条班级记录的ID拿到这一条旧的数据
	var class db.PartyClass
	if err := db.Driver.GetOne("party_class",cid,&class);err != nil{
		//这里的报错信息使用方法是:包名.类名
		panic(classException.ClassNotFount())
	}
	//todo 创建者是否当前   done
	//拿到 登陆者的ID
	accountId := auth.AccountModel().Id
	//从class这条记录中拿到 创建这条记录的老师Id,即AccountId
	if accountId != class.AccountId && !auth.IsAdmin(){
		//当前登陆者不能修改 别人的创建的班级
		panic(classException.IllegalModify())
	}

	//前端发来的 请求体
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(&class)
	//解释：params.Diff() 就是自己写的方法，如果前端传过来的请求体有这个字段，就修改，如果没有就从原来的这条记录拿。所以就不用if params.Has()
	//修改对应的数据
//
	if params.Has("name"){
		class.Name = params.Str("name","班级名称")
	}
	if params.Has("introduce"){
		class.Introduce = params.Str("introduce","简介")
	}
	if params.Has("place"){
		class.Place = params.Str("place","地点")
	}
	//int -> int64时间戳
	if params.Has("startTime"){
		class.StartTime = int64(params.Int("startTime","开始时间"))
	}
	if params.Has("endTime"){
		class.EndTime = int64(params.Int("endTime","结束时间"))
	}
	if params.Has("comment"){
		class.Comment = params.Str("comment","备注")
	}
	//保存回去db
	db.Driver.Save(&class)
	//response 告诉前端 ok
	ctx.JSON(iris.Map{
		"id": class.Id,
	})
}

//Delete     ---  Delete     删除  cid int
func DeletePartyClass(ctx iris.Context,auth authbase.AuthAuthorization,cid int){
	auth.CheckLogin()

	var class db.PartyClass
	//从db or cache get one data
	if err := db.Driver.GetOne("party_class",cid,&class);err == nil{
		//成功拿到这条记录
		//todo 判断登陆者是不是创建者   done
		if auth.AccountModel().Id != class.AccountId && !auth.IsAdmin(){
			panic(classException.IllegalDelete())
		}else{
			//db delete
			db.Driver.Delete(class)
		}
	}
	//response
	ctx.JSON(iris.Map{
		"id":cid,
	})
}

////参考car.go购物车的MDelete
//func MDeleteClass(ctx iris.Context, auth authbase.AuthAuthorization) {
//	auth.CheckLogin()
//
//	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
//	ids := params.List("ids","课程列表")
//	if err := db.Driver.Debug().Exec("delete from class1 where id in (?)",ids).Error; err != nil {
//
//		fmt.Println(err)
//		panic(classException.ClassNotFount())
//	}
//
//	ctx.JSON(iris.Map{
//		"id": ids,
//	})
//}

//获取很多条记录  List --- get    获取n个课程的  id，名称，等 很少个字段
func ListPartyClasses(ctx iris.Context,auth authbase.AuthAuthorization){
	//自定义的list，很少个但是很关键的字段
	//(比如用户进来首页面的时候，需要加载很多商品，然后前端就可以用Mget()方法从这个list里面根据Name,id来获取更加详细的商品信息)
	//这里定义的struct是table的字段的集合的子集  感觉像Model
	auth.CheckLogin()
	//list接口一般都是拿id和createTime
	var lists []struct {
		Id         int   `json:"id"`
		CreateTime int64 `json:"create_time"`
	}
	//多少条记录
	var count int
	//specify the table on which you would like to run db operations
	table := db.Driver.Table("party_class")

	//这些都是url路径参数，不填就default。
	//一页多少条记录
	limit := ctx.URLParamIntDefault("limit", 10)
	//分页
	page := ctx.URLParamIntDefault("page", 1)
	//todo 学生看到自己加入的班级
	//todo 管理员可以看到全部，老师可以看到自己创建的课程
	//教师       因为学生不可能创建班级的，所以学生的accountId根本不可能存在于Class1的accountId
	if !auth.IsAdmin() {
		table = table.Where("account_id = ?", auth.AccountModel().Id)
	}
	//管理员
	if author := ctx.URLParamIntDefault("author_id", 0); author != 0 && auth.IsAdmin() {
		table = table.Where("account_id = ?", author)
	}
	////todo 学生只能看到自己的班级
	//班级-学生表   n-n关系   Id   studentId   classId
	//sctable := db.Driver.Table("student_class")
	//////这就是一个select * from where
	//sctable = sctable.Where("account_id = ?",auth.AccountModel().Id)
	//
	//var classIds []int
	////班级-学生的那张表 -》根据学生id 去找到学生参加的班级，return ids
	//sctable.Select("class_id").Find(&classIds)
	//
	//logUtils.Println(classIds)
	//
	//
	//
	//
	//table = table.Where("id in (?)",classIds)
                                                                     //这里也小写
	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, create_time").Find(&lists)
	//向前端返回 lists
	ctx.JSON(iris.Map{
		"classes": lists,
		"total": count,         //总共多少条
		"limit": limit,          //一页多少条
		"page":  page,        //当前页
	})
}

//todo  add some words    done
var x = []string{
	"Id","Name","AccountId","PartyCourseId","Introduce","Code","Place","StartTime","EndTime","Comment","CreateTime","UpdateTime","TeacherName",
}

//Mget --- post    根据前端给来的 id 数组，进行获取更详细的goods信息
func MgetPartyClass(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	classes := db.Driver.GetMany("party_class", ids, db.PartyClass{})
	//for   goods  ->  data
	for _,class  := range classes {
		func(data *[]interface{}) {
			//这里对应的是数据库表的字段名 Select:选择“指定查询时要从数据库检索的字段”，默认情况下，将选择所有字段；创建/更新时，指定要保存到数据库的字段
			*data = append(*data, paramsUtils.ModelToDict(class,x))
			defer func() {
				recover()
			}()
		}(&data)
	}
	//返回data
	ctx.JSON(data)
}


