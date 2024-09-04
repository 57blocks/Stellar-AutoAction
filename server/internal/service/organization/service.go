package organization

import (
	"context"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/restyx"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/organization"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type (
	Service interface {
		Organization(c context.Context) (*model.Organization, error)
		OrgRoleKey(c context.Context) (*dto.RespRelatedRoleKey, error)
		OrgSecret(c context.Context) (string, error)

		KeepRootSessionAlive(c context.Context) error
	}
	conductor struct{}
)

var (
	Conductor Service
)

func init() {
	if Conductor == nil {
		Conductor = &conductor{}
	}
}

func (cd conductor) Organization(c context.Context) (*model.Organization, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errors.New("convert context.Context to gin.Context failed")
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	org := new(model.Organization)
	if err := db.Conn(c).Table(org.TableName()).
		Where(map[string]interface{}{
			"name": jwtOrg,
		}).
		First(org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("none organization found")
		}

		return nil, errors.Wrap(err, "db error when find organization by name")
	}

	return org, nil
}

func (cd conductor) OrgRoleKey(c context.Context) (*dto.RespRelatedRoleKey, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errors.New("convert context.Context to gin.Context failed")
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	orkList := make([]*model.CSOrgRoleKey, 0)
	if err := db.Conn(c).
		Table(model.TabNameAbbrOrgRoleKey()).
		Joins("LEFT JOIN organization AS o ON o.id = ork.organization_id").
		Where(map[string]interface{}{
			"o.name": jwtOrg,
		}).
		Find(&orkList).Error; err != nil {
		return nil, errors.Wrap(err, "db error when query organization key pairs")
	}

	if len(orkList) == 0 {
		return nil, errors.New("none organization related key pairs found")
	}

	csRoleKeyList := make([]dto.RespCSRoleKey, 0, len(orkList))
	for _, ork := range orkList {
		csRoleKeyList = append(csRoleKeyList, dto.RespCSRoleKey{
			CSRoleID: ork.CSRoleID,
			CSKeyID:  ork.CSKeyID,
			CSScopes: ork.CSScopes,
		})
	}

	return &dto.RespRelatedRoleKey{
		CSRoleKeys: csRoleKeyList,
	}, nil
}

func (cd conductor) OrgSecret(c context.Context) (string, error) {
	// TODO: get secrets from Amazon Secrets Manager

	return "", nil
}

// KeepRootSessionAlive
// Login through `cs` manually at the first time, then save the required info to DB,
// Here below will keep the root session alive by calling the refresh endpoint
// with a fixed interval.
//
// The interval is based on the `--refresh-lifetime` in root session initialization:
//
//	`cs login -s google --session-lifetime 31536000 --auth-lifetime 86400 --refresh-lifetime 86400`
//
// The longest session lifetime is 31536000 seconds, which is about 1 year.
//
// Once the root session is expired, login through `cs` again and sync the info to DB.
func (cd conductor) KeepRootSessionAlive(c context.Context) error {
	rootSessions := make([]*model.OrgRootSession, 0)
	if err := db.Conn(c).Table(model.TabNameOrgSession()).
		Where(map[string]interface{}{}).
		Find(rootSessions).Error; err != nil {
		return errors.Wrap(err, "db error when find organization root session")
	}

	if len(rootSessions) == 0 {
		return errors.New("none organization root sessions found")
	}

	results := make([]dto.RespRootSessionRefresh, 0, len(rootSessions))
	for _, session := range rootSessions {
		response, err := restyx.Client.R().
			EnableTrace().
			SetHeaders(map[string]string{
				"Content-Type":  "application/json",
				"accept":        "application/json",
				"Authorization": session.Token,
			}).
			SetBody(map[string]interface{}{
				"epoch_num":   session.Epoch,
				"epoch_token": session.EpochToken,
				"other_token": session.EpochRefreshToken,
			}).
			Post(fmt.Sprintf(""))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("resty error: %s\n", err.Error()))
		}
		fmt.Println(response.String())

	}

	return nil
}
