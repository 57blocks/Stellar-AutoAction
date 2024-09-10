package model

// CubeSignerKey is a struct that represents the required information
// that used for signing.
type CubeSignerKey struct {
	ICU
	AccountID uint64  `json:"account_id"`
	Key       string  `json:"key"`
	Scopes    StrList `json:"scopes" gorm:"type:text[]"`
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
