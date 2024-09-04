package model

type Organization struct {
	ICU
	Name             string `json:"name"`
	CSOrganizationID string `json:"cs_organization_id"`
	Description      string `json:"description"`
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

// CSOrgRoleKey is a struct that represents the organization key pairs.
// In model, it's a 1:1 mapping between role and key, which aims at
// minimizing the scopes of a role. Like a role named `Signer` could
// only have one key to have the only access with sign.
type CSOrgRoleKey struct {
	ICU
	OrganizationID uint64  `json:"organization_id"`
	CSRoleID       string  `json:"cs_role_id"`
	CSKeyID        string  `json:"cs_key_id"`
	CSScopes       StrList `json:"cs_scopes" gorm:"column:cs_scopes;type:text[]"`
}

func (o *CSOrgRoleKey) TableName() string {
	return "organization_role_key"
}

func (o *CSOrgRoleKey) TableNameWithAbbr() string {
	return "organization_role_key AS ork"
}

func TabNameOrgRoleKey() string {
	return (&CSOrgRoleKey{}).TableName()
}

func TabNameOrgRoleKeyAbbr() string {
	return (&CSOrgRoleKey{}).TableNameWithAbbr()
}
