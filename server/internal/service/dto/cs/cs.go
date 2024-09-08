package cs

import "github.com/57blocks/auto-action/server/internal/model"

// Role model dto
type Role struct {
	ID uint64 `json:"id"`
	//OrganizationID uint64 `json:"organization_id"`
	//AccountID      uint64 `json:"account_id"`
	Role string `json:"role"`
	//Keys []RespToSignKey `json:"keys"`
}

//Key struct {
//	//ID     uint64        `json:"id"`
//	//RoleID uint64        `json:"role_id"`
//	Key    string        `json:"key"`
//	Scopes model.StrList `json:"scopes" gorm:"type:text[]"`
//}

type (
	RespToSign struct {
		_            struct{}
		Organization string          `json:"organization"`
		Role         string          `json:"role"`
		Keys         []RespToSignKey `json:"keys"`
	}

	RespToSignKey struct {
		_      struct{}
		ID     uint64        `json:"id"`
		Key    string        `json:"key"`
		Scopes model.StrList `json:"scopes"`
	}
)

type (
	ReqToSign struct {
		Organization string `json:"organization"`
		Account      string `json:"account"`
	}
)
