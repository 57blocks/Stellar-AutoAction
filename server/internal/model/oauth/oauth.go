package oauth

import (
	"time"

	"github.com/57blocks/auto-action/server/internal/model"
)

type Organization struct {
	model.ICU
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

func TabNamOrg() string {
	return (&Organization{}).TableName()
}

func TabNamOrgAbbr() string {
	return (&Organization{}).TableNameWithAbbr()
}

type User struct {
	model.ICU
	Account        string `json:"account"`
	Password       string `json:"password"`
	UserKey        string `json:"user_key"`
	Description    string `json:"description"`
	OrganizationId int32  `json:"organization_id"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) TableNameAbbr() string {
	return "\"user\" AS u"
}

func TabNamUser() string {
	return (&User{}).TableName()
}

func TabNamUserAbbr() string {
	return (&User{}).TableNameAbbr()
}

type Token struct {
	model.ICU
	Access         string    `json:"access"`
	Refresh        string    `json:"refresh"`
	UserId         uint64    `json:"user_id"`
	AccessExpires  time.Time `json:"access_expires"`
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
