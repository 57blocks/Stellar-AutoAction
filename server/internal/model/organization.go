package model

type Organization struct {
	ICU
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (o *Organization) TableName() string {
	return "principal_organization"
}

func (o *Organization) TableNameWithAbbr() string {
	return "principal_organization AS po"
}
