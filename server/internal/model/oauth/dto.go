package oauth

type (
	ReqID struct {
		ID uint64 `json:"id"`
	}
	ReqName struct {
		Name string `json:"name"`
	}
)

// User model representations in request
type (
	ReqOrgAcn struct {
		OrgName string `json:"org_name"`
		AcnName string `json:"acn_name"`
	}

	RespUser struct {
		ID             uint64   `json:"id"`
		Account        string   `json:"account"`
		Password       string   `json:"-"`
		Description    string   `json:"description"`
		OrganizationId int32    `json:"-"`
		Organization   *RespOrg `json:"organization,omitempty" gorm:"foreignKey:organization_id"`
	}
)

// RespOrg organization related dto
type RespOrg struct {
	ID            uint64 `json:"-"`
	Name          string `json:"name"`
	CubeSignerOrg string `json:"cube_signer_org"`
	Description   string `json:"description"`
}

// Token related dto
type ()
