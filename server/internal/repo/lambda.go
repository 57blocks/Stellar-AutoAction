package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/util"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"gorm.io/gorm"
)

//go:generate mockgen -destination ./lambda_mock.go -package repo -source lambda.go Lambda
type (
	Lambda interface {
		FindByNameOrARN(c context.Context, input string) (*dto.RespInfo, error)
		LambdaInfo(c context.Context, req *dto.ReqInfo) (*dto.RespInfo, error)
		PersistRegResult(c context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
		FindByAccountId(c context.Context, accountId uint64) ([]*dto.RespInfo, error)
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

func (l *lambda) FindByNameOrARN(c context.Context, input string) (*dto.RespInfo, error) {
	lamb := new(dto.RespInfo)

	if err := l.Instance.Conn(c).Table("lambda").
		Where(map[string]interface{}{
			"function_arn": input,
		}).
		Or(map[string]interface{}{
			"function_name": util.GenLambdaFuncName(c, input),
		}).
		First(lamb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound(fmt.Sprintf("none lambda found by: %s", input))
		}

		return nil, errorx.Internal(fmt.Sprintf("failed to query lambda: %s, err: %s", input, err.Error()))
	}

	return lamb, nil
}

func (l *lambda) LambdaInfo(c context.Context, req *dto.ReqInfo) (*dto.RespInfo, error) {
	resp := new(dto.RespInfo)

	if err := l.Instance.Conn(c).Table("lambda").
		Preload("Schedulers", func(db *gorm.DB) *gorm.DB {
			return db.Table("lambda_scheduler")
		}).
		Where(map[string]interface{}{
			"function_arn": req.Lambda,
		}).
		Or(map[string]interface{}{
			"function_name": util.GenLambdaFuncName(c, req.Lambda),
		}).
		First(resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound(fmt.Sprintf("none lambda found by: %s", req.Lambda))
		}

		return nil, errorx.Internal(fmt.Sprintf("failed to query lambda: %s, err: %s", req.Lambda, err.Error()))
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

func (l *lambda) FindByAccountId(c context.Context, accountId uint64) ([]*dto.RespInfo, error) {
	resp := make([]*dto.RespInfo, 0)

	if err := l.Instance.Conn(c).Table("lambda").
		Where(map[string]interface{}{
			"account_id": accountId,
		}).Find(&resp).Error; err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to query lambda by account, err: %s", err.Error()))
	}

	return resp, nil
}
