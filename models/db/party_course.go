package db

// 党课表
type PartyCourse struct {
	Id int `gorm:"primary_key" json:"id"`

	// 课程名称
	Name string `json:"party_course_name"`

	AccountId int `json:"account_id" gorm:"not null;index"`

	// 课程简介
	CourseBrief string `json:"party_course_brief" gorm:"type:text"`

	// 课程封面
	CourseCover string `json:"party_course_cover"`

	// 作业简介
	CourseWork string`json:"party_course_work" gorm:"type:text"`

	// 商品ID
	GoodsId int `json:"good_id"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}

// 资料表
type PartyCourseData struct {
	Id int `gorm:"primary_key" json:"id"`

	// 党课ID
	PartyCourseId int `json:"party_course_id"`

	// 课程ppt
	CoursePpt string `json:"party_course_ppt"`

	// 课程视频
	CourseVideo string `json:"party_course_video"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}

// 课程轮播图
type CoursePicture struct {
	Id int `gorm:"primary_key" json:"id"`

	// 课程列表
	CourseList string `json:"course_list"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}
