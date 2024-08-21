package oauth

import (
	"context"
	"time"

	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/db"
	"github.com/57blocks/auto-action/server/internal/pkg/jwtx"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/oauth"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	Service interface {
		Login(c context.Context, request dto.Request) (*dto.Response, error)
		Logout(c context.Context, req dto.Request)
	}
	ServiceConductor struct{}
)

var Conductor Service

func init() {
	if Conductor == nil {
		Conductor = &ServiceConductor{}
	}
}

func (sc *ServiceConductor) Login(c context.Context, req dto.Request) (*dto.Response, error) {
	user := new(model.User)
	if err := db.Conn(c).Table(user.TableNameWithAbbr()).
		//Select("*").
		Joins("LEFT JOIN principal_organization AS po ON pu.organization_id = po.id").
		Where(map[string]interface{}{
			"pu.account": req.Account,
			"po.name":    req.Organization,
		}).
		First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user/organization not found")
		}
		return nil, errors.New(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), req.Password); err != nil {
		return nil, errors.New("password not match")
	}

	// tokens assignment
	now := time.Now().UTC()
	tokenExp := now.AddDate(0, 1, 0)
	refreshExp := now.AddDate(0, 3, 0)

	tokens, err := jwtx.Assign(tokenExp, refreshExp)
	if err != nil {
		return nil, err
	}

	resp := dto.BuildResp(
		dto.WithAccount(req.Account),
		dto.WithOrganization(req.Organization),
		dto.WithBound(dto.BuildBound(
			dto.WithBoundName(req.Environment),
			dto.WithBoundEndPoint(viper.GetString("bound.endpoint")),
		)),
		dto.WithTokens(tokens),
	)

	// save tokens association
	token := &model.Token{
		Access:         resp.Token,
		Refresh:        resp.Refresh,
		UserId:         user.ID,
		AccessExpires:  tokenExp,
		RefreshExpires: refreshExp,
	}
	if err := db.Conn(c).Create(token).Error; err != nil {
		return nil, errors.New(err.Error())
	}

	return resp, nil
}

func (sc *ServiceConductor) Logout(c context.Context, req dto.Request) {

}
