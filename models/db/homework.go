package db

type Homework struct {
	Id int `json:"id" gorm:"primary_key"`
	//上传者ID   学生ID
	UpperId int `json:"upper_id"`
	//班级id
	ClassId int `json:"class_id"`
	//课程id
	CourseId int `json:"course_id"`

	Content string `json:"content"`
	//图片
	Picture string `json:"picture"`
	//视频
	Video string `json:"video"`
	//创建时间
	CreateTime int64 `json:"create_time"`
	//更新时间
	UpdateTime int64 `json:"update_time"`
}
