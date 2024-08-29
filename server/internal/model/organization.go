package model

type Organization struct {
	ICU
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (o *Organization) TableName() string {
	return "organization"
}

func (o *Organization) TableNameWithAbbr() string {
	return "organization AS o"
}

type OrganizationKeyPairs struct {
	ICU
	OrganizationID string `json:"organization_id"`
	PublicKey      string `json:"public_key"`
	PrivateKey     string `json:"private_key"`
}

func (o *OrganizationKeyPairs) TableName() string {
	return "principal_org_key_pairs"
}

func (o *OrganizationKeyPairs) TableNameWithAbbr() string {
	return "organization_key_pairs AS okp"
}
