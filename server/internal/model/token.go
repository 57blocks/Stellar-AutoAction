package model

import (
	"time"
)

type Token struct {
	ICU
	Access         string    `json:"access"`
	Refresh        string    `json:"refresh"`
	UserId         int64     `json:"user_id"`
	AccessExpires  time.Time `json:"access_expires"`
	RefreshExpires time.Time `json:"refresh_expires"`
}

func (t *Token) TableName() string {
	return "principal_token"
}

func (t *Token) TableNameWithAbbr() string {
	return "principal_token AS pt"
}
