package model

// CubeSignerRole is a struct that represents the organization roles in
// the CubeSigner.
// Each account will host one Role only, with it named by the account.
type CubeSignerRole struct {
	ICU
	OrganizationID uint64 `json:"organization_id"`
	AccountID      uint64 `json:"account_id"`
	Role           string `json:"role"`
}

func (o *CubeSignerRole) TableName() string {
	return "cube_signer_role"
}

func (o *CubeSignerRole) TableNameWithAbbr() string {
	return "cube_signer_role AS csr"
}

func TabNamCSRole() string {
	return (&CubeSignerRole{}).TableName()
}

func TabNamCSRoleAbbr() string {
	return (&CubeSignerRole{}).TableNameWithAbbr()
}

// CubeSignerKey is a struct that represents the keys those are bound
// to a Role in the CubeSigner.
type CubeSignerKey struct {
	ICU
	RoleID uint64  `json:"role_id"`
	Key    string  `json:"key"`
	Scopes StrList `json:"scopes" gorm:"type:text[]"`
}

func (o *CubeSignerKey) TableName() string {
	return "cube_signer_key"
}

func (o *CubeSignerKey) TableNameWithAbbr() string {
	return "cube_signer_key AS csk"
}

func TabNameCSKey() string {
	return (&CubeSignerKey{}).TableName()
}

func TabNameCSKeyAbbr() string {
	return (&CubeSignerKey{}).TableNameWithAbbr()
}
