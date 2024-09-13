package dto

import (
	"time"

	"github.com/57blocks/auto-action/server/internal/third-party/jwtx"
)

// general representation
type (
	ReqID struct {
		ID uint64 `json:"id"`
	}
	ReqName struct {
		Name string `json:"name"`
	}
)

// Login related dto
type (
	ReqLogin struct {
		_            struct{}
		Account      string `json:"account"`
		Organization string `json:"organization"`
		Password     []byte `json:"password"`
		Environment  string `json:"environment"`
	}

	RespCredential struct {
		_              struct{}
		Account        string `json:"account" toml:"account"`
		Organization   string `json:"organization" toml:"organization"`
		Environment    string `json:"environment" toml:"environment"`
		jwtx.TokenPair `json:"tokens" toml:"tokens"`
	}
	RespCredOpt func(cred *RespCredential)
)

// BuildRespCred build the RespCredential of Login
func BuildRespCred(opts ...RespCredOpt) *RespCredential {
	cred := new(RespCredential)

	for _, opt := range opts {
		opt(cred)
	}

	return cred
}

func WithAccount(account string) RespCredOpt {
	return func(resp *RespCredential) {
		resp.Account = account
	}
}

func WithOrganization(organization string) RespCredOpt {
	return func(resp *RespCredential) {
		resp.Organization = organization
	}
}

func WithEnvironment(environment string) RespCredOpt {
	return func(resp *RespCredential) {
		resp.Environment = environment
	}
}

func WithTokenPair(tp jwtx.TokenPair) RespCredOpt {
	return func(resp *RespCredential) {
		resp.TokenPair = tp
	}
}

// RespLogout related dto
type RespLogout struct{}

// User model representations in request
type (
	ReqOrgAcn struct {
		OrgName string `json:"org_name"`
		AcnName string `json:"acn_name"`
	}

	RespUser struct {
		ID             uint64   `json:"-"`
		Account        string   `json:"account"`
		Password       string   `json:"-"`
		Description    string   `json:"-"`
		CubeSignerUser string   `json:"cube_signer_user"`
		OrganizationId int32    `json:"-"`
		Organization   *RespOrg `json:"organization,omitempty" gorm:"foreignKey:organization_id"`
	}
)

// RespOrg organization related dto
type RespOrg struct {
	ID          uint64 `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type (
	RespToken struct {
		ID             uint64     `json:"-"`
		Access         string     `json:"access"`
		Refresh        string     `json:"refresh"`
		UserId         uint64     `json:"user_id"`
		AccessExpires  time.Time  `json:"access_expires"`
		RefreshExpires time.Time  `json:"refresh_expires"`
		UpdatedAt      *time.Time `json:"updated_at"`
	}
)
