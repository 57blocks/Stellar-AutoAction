package cs

import (
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/model/oauth"
)

// CubeSigner
type (
	ReqToSign struct {
		Organization string `json:"organization"`
		Account      string `json:"account"`
	}

	RespToSign struct {
		_              struct{}
		ID             uint64         `json:"-"`
		OrganizationID uint64         `json:"-"`
		Organization   oauth.RespOrg  `json:"organization" gorm:"foreignKey:organization_id;references:id"`
		AccountID      uint64         `json:"-"`
		Account        oauth.RespUser `json:"account" gorm:"foreignKey:account_id;references:id"`
		Role           string         `json:"role"`
		Keys           []RespCSKey    `json:"keys" gorm:"foreignKey:role_id;references:id"`
	}

	RespCSKey struct {
		_      struct{}
		ID     uint64        `json:"-"`
		RoleID uint64        `json:"-"`
		Key    string        `json:"key"`
		Scopes model.StrList `json:"scopes" gorm:"type:text[]"`
	}
)
