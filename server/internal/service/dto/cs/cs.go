package cs

import "github.com/57blocks/auto-action/server/internal/model"

type (
	RespToSign struct {
		_            struct{}
		ID           uint64          `json:"-"`
		Organization string          `json:"organization"`
		Role         string          `json:"role"`
		Keys         []RespToSignKey `json:"keys" gorm:"foreignKey:role_id;references:id"`
	}

	RespToSignKey struct {
		_      struct{}
		ID     uint64        `json:"-"`
		RoleID uint64        `json:"-"`
		Key    string        `json:"key"`
		Scopes model.StrList `json:"scopes" gorm:"type:text[]"`
	}
)

type (
	ReqToSign struct {
		Organization string `json:"organization"`
		Account      string `json:"account"`
	}
)
