package dto

import "github.com/57blocks/auto-action/server/internal/model"

// CubeSigner
type (
	ReqToSign struct {
		Organization string `json:"organization"`
		Account      string `json:"account"`
		From         string `json:"from"`
	}

	RespCSKey struct {
		_            struct{}
		ID           uint64        `json:"-"`
		AccountID    uint64        `json:"-"`
		Account      RespUser      `json:"account,omitempty" gorm:"foreignKey:account_id;references:id"`
		Organization string        `json:"organization"`
		Role         string        `json:"role"`
		Key          string        `json:"key"`
		Scopes       model.StrList `json:"scopes" gorm:"type:text[]"`
	}

	AddCsKeyResponse struct {
		Keys []struct {
			KeyID string `json:"key_id"`
		} `json:"keys"`
	}
)
