package db

type SignIn struct{
	ID int `grom:"primary_key" json:"id"`

	AccountID int `json:"account_id" grom:"not null;index"`
	
	Date string `json:"date"`

	CreateTime int64 `json:"create_time"`
}
