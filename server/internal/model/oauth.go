package model

import (
	"time"
)

type Organization struct {
	ICU
	Name          string `json:"name"`
	CubeSignerOrg string `json:"cube_signer_org"`
	Description   string `json:"description"`
}

func (o *Organization) TableName() string {
	return "organization"
}

func (o *Organization) TableNameWithAbbr() string {
	return "organization AS o"
}

func TabNameOrg() string {
	return (&Organization{}).TableName()
}

func TabNameOrgAbbr() string {
	return (&Organization{}).TableNameWithAbbr()
}

type User struct {
	ICU
	Account        string `json:"account"`
	Password       string `json:"password"`
	Description    string `json:"description"`
	OrganizationId int32  `json:"organization_id"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) TableNameAbbr() string {
	return "\"user\" AS u"
}

func TabNameUser() string {
	return (&User{}).TableName()
}

func TabNameUserAbbr() string {
	return (&User{}).TableNameAbbr()
}

type Token struct {
	ICU
	UserId         uint64    `json:"user_id"`
	Access         string    `json:"access"`
	AccessID       string    `json:"access_id"`
	AccessExpires  time.Time `json:"access_expires"`
	Refresh        string    `json:"refresh"`
	RefreshID      string    `json:"refresh_id"`
	RefreshExpires time.Time `json:"refresh_expires"`
}

func (t *Token) TableName() string {
	return "token"
}

func (t *Token) TableNameWithAbbr() string {
	return "token AS t"
}

func TabNameToken() string {
	return (&Token{}).TableName()
}

func TabNameTokenAbbr() string {
	return (&Token{}).TableNameWithAbbr()
}
