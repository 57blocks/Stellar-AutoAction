package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/util"

	"gorm.io/gorm"
)

//go:generate mockgen -destination ./lambda_mock.go -package repo -source lambda.go Lambda
type (
	Lambda interface {
		Info(c context.Context, req *dto.ReqInfo) (*dto.RespInfo, error)
	}
	lambda struct{}
)

var CDLambda Lambda

func init() {
	if CDLambda == nil {
	}
	CDLambda = &lambda{}
}
func (l *lambda) Info(c context.Context, req *dto.ReqInfo) (*dto.RespInfo, error) {
	resp := new(dto.RespInfo)

	if err := db.Conn(c).Table("lambda").
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
