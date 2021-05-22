package db

type OrderInfo struct {
	ID int `json:"id" gorm:"primary_key"`

	Number string `json:"number" gorm:"not null"`

	GoodsID int `json:"goods_id" gorm:"not null;index"`

	AccountID int 	`json:"account_id" gorm:"not null;index"`

	Total int `json:"total" gorm:"default:1"`

	TotalPrice int `json:"total_price"`

	CreateTime int64 `json:"create_time"`

	UpdateTime int64 `json:"update_time"`
}
