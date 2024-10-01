package repo

import (
	"database/sql"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMain(m *testing.M) {
	testConfig := config.Configuration{
		Log: config.Log{
			Level:    "debug",
			Encoding: "json",
		},
	}
	logx.Setup(&testConfig)

	os.Exit(m.Run())
}

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       sqldb,
		DriverName: "postgres",
	})
	gormdb, _ := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	return sqldb, gormdb, mock
}

func TestFindCSKeySuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testKey := "testKey"
	testAccountID := uint64(123)
	rows := sqlmock.NewRows([]string{"Key", "AccountID"}).AddRow(testKey, testAccountID)
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	ctx := new(gin.Context)
	repo := &cubeSigner{
		Instance: &db.Instance{DB: gormdb},
	}
	csKey, err := repo.FindCSKey(ctx, testKey, testAccountID)

	assert.NoError(t, err)
	assert.Equal(t, testKey, csKey.Key)
	assert.Equal(t, testAccountID, csKey.AccountID)
}

func TestFindCSKeyNotFound(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	rows := sqlmock.NewRows([]string{"Key", "AccountID"})
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	ctx := new(gin.Context)
	repo := &cubeSigner{
		Instance: &db.Instance{DB: gormdb},
	}
	csKey, err := repo.FindCSKey(ctx, "testKey", 1)

	assert.Error(t, err)
	assert.Equal(t, errorx.NotFound("cube signer key not found"), err)
	assert.Nil(t, csKey)
}

func TestSyncCSKeySuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testKey := "testKey"
	testAccountID := uint64(123)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "cube_signer_key"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), testAccountID, testKey, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	ctx := new(gin.Context)
	repo := &cubeSigner{
		Instance: &db.Instance{DB: gormdb},
	}
	err := repo.SyncCSKey(ctx, &model.CubeSignerKey{Key: testKey, AccountID: testAccountID, Scopes: []string{"testScopes"}})

	assert.NoError(t, err)
}

func TestSyncCSKeyError(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testKey := "testKey"
	testAccountID := uint64(123)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "cube_signer_key"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), testAccountID, testKey, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctx := new(gin.Context)
	repo := &cubeSigner{
		Instance: &db.Instance{DB: gormdb},
	}
	err := repo.SyncCSKey(ctx, &model.CubeSignerKey{Key: testKey, AccountID: testAccountID, Scopes: []string{"testScopes"}})

	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("error"), err)
}

func TestDeleteCSKeySuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testKey := "testKey"
	testAccountID := uint64(123)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "cube_signer_key"`)).
		WithArgs(testKey, testAccountID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	ctx := new(gin.Context)
	repo := &cubeSigner{
		Instance: &db.Instance{DB: gormdb},
	}
	err := repo.DeleteCSKey(ctx, testKey, testAccountID)

	assert.NoError(t, err)
}

func TestDeleteCSKeyError(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testKey := "testKey"
	testAccountID := uint64(123)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "cube_signer_key"`)).
		WithArgs(testKey, testAccountID).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctx := new(gin.Context)
	repo := &cubeSigner{
		Instance: &db.Instance{DB: gormdb},
	}
	err := repo.DeleteCSKey(ctx, testKey, testAccountID)

	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("error"), err)
}

func TestFindCSKeysByAccountSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testAccountID := uint64(123)
	rows := sqlmock.NewRows([]string{"Key", "AccountID"}).
		AddRow("testKey1", testAccountID).
		AddRow("testKey2", testAccountID)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(rows)

	ctx := new(gin.Context)
	repo := &cubeSigner{
		Instance: &db.Instance{DB: gormdb},
	}
	csKeys, err := repo.FindCSKeysByAccount(ctx, testAccountID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(csKeys))
	assert.Equal(t, "testKey1", csKeys[0].Key)
	assert.Equal(t, "testKey2", csKeys[1].Key)
}

func TestFindCSKeysByAccountError(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	testAccountID := uint64(123)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnError(errors.New("error"))

	ctx := new(gin.Context)
	repo := &cubeSigner{
		Instance: &db.Instance{DB: gormdb},
	}
	csKeys, err := repo.FindCSKeysByAccount(ctx, testAccountID)

	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("error"), err)
	assert.Nil(t, csKeys)
}
