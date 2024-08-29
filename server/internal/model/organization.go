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

type OrgKeyPairs struct {
	ICU
	OrganizationID string `json:"organization_id"`
	PublicKey      string `json:"public_key"`
	PrivateKey     string `json:"private_key"`
}

func (o *OrgKeyPairs) TableName() string {
	return "principal_org_key_pairs"
}

func (o *OrgKeyPairs) TableNameWithAbbr() string {
	return "principal_org_key_pairs AS pokp"
}
