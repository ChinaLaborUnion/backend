package db

//资讯
type NewsInfo struct {
	Id int `json:"id" gorm:"primary_key"`
	//标题
	Title string `json:"title"`
	//简介
	Introduction string `json:"introduction" gorm:"type:text"`
	//内容
	Content string `json:"content" gorm:"type:text"`
	//标签Id
	NewsLabelId int `json:"news_label_id" gorm:"index;not null"`
	//是否发布
	IsPublish bool `json:"is_publish"`
	//封面图片
	Picture string `json:"pictures"`

	CreateTime int64 `json:"create_time"`
	UpdateTime int64 `json:"update_time"`
}

//资讯标签
type NewsLabel struct {
	Id int `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	CreateTime int64 `json:"create_time"`
	UpdateTime int64 `json:"update_time"`
}