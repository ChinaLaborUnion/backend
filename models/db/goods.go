package db

type GoodsInfo struct {
	ID int `json:"id" gorm:"primary_key"`

	Cover string `json:"cover" gorm:"not null"`

	Name string  `json:"name" gorm:"not null"`

	Brief string `json:"brief"`

	Price int `json:"price" gorm:"not null"`

	People int `json:"people" gorm:"default:0"`

	Pictures string  `json:"pictures"`

	CreateTime int64 `json:"create_time"`

	UpdateTime int64 `json:"update_time"`
}


