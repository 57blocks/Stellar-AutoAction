package repo

import (
	"errors"
	"testing"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFindUserByAcnSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testAcn := "testAcn"
	rows := sqlmock.NewRows([]string{"Account", "OrganizationId"}).
		AddRow(testAcn, 1)
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	user, err := repo.FindUserByAcn(ctx, testAcn)

	assert.NoError(t, err)
	assert.Equal(t, testAcn, user.Account)
	assert.Equal(t, int32(1), user.OrganizationId)
}

func TestFindUserByAcnNotFound(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testAcn := "testAcn"
	mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	user, err := repo.FindUserByAcn(ctx, testAcn)

	assert.Error(t, err)
	assert.Equal(t, "user/organization not found", err.Error())
	assert.Nil(t, user)
}

func TestFindUserByOrgAcnSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testAcn := "testAcn"
	testOrg := "testOrg"
	rows := sqlmock.NewRows([]string{"Account", "OrganizationId"}).
		AddRow(testAcn, 1)
	expectedSQL := `SELECT .+ FROM "user" AS u LEFT JOIN organization AS o ON u.organization_id = o.id`
	mock.ExpectQuery(expectedSQL).WithArgs(testAcn, testOrg, 1).
		WillReturnRows(rows)

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	user, err := repo.FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		AcnName: testAcn,
		OrgName: testOrg,
	})

	assert.NoError(t, err)
	assert.Equal(t, testAcn, user.Account)
	assert.Equal(t, int32(1), user.OrganizationId)
}

func TestFindUserByOrgAcnNotFound(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testAcn := "testAcn"
	testOrg := "testOrg"
	mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	user, err := repo.FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		AcnName: testAcn,
		OrgName: testOrg,
	})

	assert.Error(t, err)
	assert.Equal(t, "user/organization not found", err.Error())
	assert.Nil(t, user)
}

func TestCreateUserSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	user := &model.User{
		Account:        "testUser",
		OrganizationId: 1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			user.Account,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			user.OrganizationId,
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	err := repo.CreateUser(ctx, user)

	assert.NoError(t, err)
}

func TestCreateUserError(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	user := &model.User{
		Account:        "testUser",
		OrganizationId: 1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			user.Account,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			user.OrganizationId,
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	err := repo.CreateUser(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, "error", err.Error())
}

func TestFindOrgByNameSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	org := &dto.RespOrg{
		ID:   1,
		Name: "testOrg",
	}

	rows := sqlmock.NewRows([]string{"Id", "Name"}).
		AddRow(org.ID, org.Name)
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	org, err := repo.FindOrgByName(ctx, org.Name)

	assert.NoError(t, err)
	assert.Equal(t, org.ID, org.ID)
	assert.Equal(t, org.Name, org.Name)
}

func TestFindOrgByNameNotFound(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	org := &dto.RespOrg{
		Name: "testOrg",
	}

	mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	org, err := repo.FindOrgByName(ctx, org.Name)

	assert.Error(t, err)
	assert.Equal(t, "organization not found", err.Error())
	assert.Nil(t, org)
}

func TestFindTokenByRefreshIDSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	refreshToken := &model.Token{
		RefreshID: "testRefreshID",
		UserId:    1,
	}

	rows := sqlmock.NewRows([]string{"refresh_id", "user_id"}).
		AddRow(refreshToken.RefreshID, refreshToken.UserId)
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	token, err := repo.FindTokenByRefreshID(ctx, refreshToken.RefreshID)

	assert.NoError(t, err)
	assert.Equal(t, refreshToken.RefreshID, token.RefreshID)
	assert.Equal(t, refreshToken.UserId, token.UserId)
}

func TestFindTokenByRefreshIDNotFound(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	refreshToken := &model.Token{
		RefreshID: "testRefreshID",
	}

	mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	token, err := repo.FindTokenByRefreshID(ctx, refreshToken.RefreshID)

	assert.Error(t, err)
	assert.Equal(t, "none refresh token found", err.Error())
	assert.Nil(t, token)
}

func TestSyncTokenSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	token := &model.Token{
		UserId: 1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	err := repo.SyncToken(ctx, token)

	assert.NoError(t, err)
}

func TestSyncTokenError(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	token := &model.Token{
		UserId: 1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO`).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	err := repo.SyncToken(ctx, token)

	assert.Error(t, err)
	assert.Equal(t, "error", err.Error())
}

func TestDeleteTokenByAccessSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	access := "testAccess"
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM`).WithArgs(access).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	err := repo.DeleteTokenByAccess(ctx, access)

	assert.NoError(t, err)
}

func TestDeleteTokenByAccessError(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	access := "testAccess"
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM`).WithArgs(access).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctx := new(gin.Context)
	repo := &oauth{
		Instance: &db.Instance{DB: gormdb},
	}

	err := repo.DeleteTokenByAccess(ctx, access)

	assert.Error(t, err)
	assert.Equal(t, "error", err.Error())
}
