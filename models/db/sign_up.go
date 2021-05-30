package db

type SignUp struct {
	Id int `grom:"primary_key" json:"id"`
	//用户ID  学生ID
	UserId int `json:"user_id" gorm:"index"`
	//班级ID
	ClassId int `json:"class_id"`
	//课程ID
	CourseId int `json:"course_id"`
	Status int16 `json:"status"`
}
