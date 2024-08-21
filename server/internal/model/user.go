package model

type User struct {
	ICU
	Account        string `json:"account"`
	Password       string `json:"password"`
	Description    string `json:"description"`
	OrganizationId int32  `json:"organization_id"`
}

func (u *User) TableName() string {
	return "principal_user"
}

func (u *User) TableNameWithAbbr() string {
	return "principal_user AS pu"
}
