package model

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
