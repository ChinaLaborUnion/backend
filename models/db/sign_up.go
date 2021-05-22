package db

type SignUp struct {
	Id int `grom:"primary_key" json:"id"`
	UserId int `json:"user_id" gorm:"index"`
	ClassId int `json:"class_id"`
	CourseId int `json:"course_id"`
	Status int16 `json:"status"`
}
