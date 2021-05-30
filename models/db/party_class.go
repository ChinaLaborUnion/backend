package db

type PartyClass struct {
	Id int `gorm:"primary_key" json:"id"`
	//授课老师ID 也是创建者
	AccountId int `json:"account_id" gorm:"not null;index"`
	//todo 党课id   done  1个党课可以创建n个班级     外键    done
	PartyCourseId int `json:"party_course_id" gorm:"not null;index"`
	//班级名称
	Name string `gorm:"not null" json:"name"`
	//班级简介
	Introduce string `json:"introduce"`
	//课程码  6位随机数字or字母
	Code string `json:"code"`
	//地点
	Place string `json:"place"`
	//开班时间
	StartTime int64 `json:"start_time"`
	//结束时间
	EndTime int64 `json:"end_time"`
	//备注
	Comment string `json:"comment"`
	//教师名字
	TeacherName string `json:"teacher_name"`
	//创建时间
	CreateTime int64 `json:"create_time"`
	//更新时间
	UpdateTime int64 `json:"update_time"`
}


