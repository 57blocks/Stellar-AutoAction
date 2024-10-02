package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/util"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"gorm.io/gorm"
)

//go:generate mockgen -destination ../testdata/lambda_mock.go -package testdata -source lambda.go Lambda
type (
	Lambda interface {
		LambdaInfo(c context.Context, acnID uint64, distinguish string) (*dto.RespInfo, error)
		PersistRegResult(c context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
		FindByAccount(c context.Context, accountId uint64) ([]*dto.RespInfo, error)
		DeleteLambdaTX(c context.Context, f func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
	}
	lambda struct {
		Instance *db.Instance
	}
)

var LambdaRepo Lambda

func NewLambda() {
	if LambdaRepo == nil {
		LambdaRepo = &lambda{
			Instance: db.Inst,
		}
	}
}

func (l *lambda) LambdaInfo(c context.Context, acnID uint64, distinguish string) (*dto.RespInfo, error) {
	resp := new(dto.RespInfo)

	if err := l.Instance.Conn(c).Table(model.TabNameLambda()).
		Preload("Scheduler", func(db *gorm.DB) *gorm.DB {
			return db.Table(model.TabNameLambdaSch())
		}).
		Where("account_id = ? and function_arn = ?", acnID, distinguish).
		Or("account_id = ? and function_name = ?", acnID, util.GenLambdaFuncName(c, distinguish)).
		First(resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound(fmt.Sprintf("none lambda found by: %s", distinguish))
		}

		return nil, errorx.Internal(fmt.Sprintf("failed to query lambda: %s, err: %s", distinguish, err.Error()))
	}

	return resp, nil
}

func (l *lambda) PersistRegResult(c context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	if err := l.Instance.Conn(c).Transaction(fc, opts...); err != nil {
		logx.Logger.ERROR(err.Error())
		return errorx.Internal("failed to persist lambda registration result")
	}

	return nil
}

func (l *lambda) FindByAccount(c context.Context, accountId uint64) ([]*dto.RespInfo, error) {
	resp := make([]*dto.RespInfo, 0)

	if err := l.Instance.Conn(c).Table(model.TabNameLambda()).
		Preload("Scheduler", func(db *gorm.DB) *gorm.DB {
			return db.Table(model.TabNameLambdaSch())
		}).
		Where("account_id = ?", accountId).
		Find(&resp).Error; err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to query lambda by account, err: %s", err.Error()))
	}

	return resp, nil
}

func (l *lambda) DeleteLambdaTX(c context.Context, f func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	if err := l.Instance.Conn(c).Transaction(f, opts...); err != nil {
		logx.Logger.ERROR(fmt.Sprintf("remove lambda and its scheduler failed, err: %s", err.Error()))

		return errorx.Internal("failed to remove lambda")
	}

	return nil
}
