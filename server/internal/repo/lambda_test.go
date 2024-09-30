package repo

import (
	"errors"
	"fmt"
	"testing"

	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/pkg/util"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFindByNameOrARNSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimIss.Str(), "test-org")
	ctx.Set(constant.ClaimSub.Str(), "test-account")
	testFuncName := "testFunc"
	testFuncArn := "testArn"
	testSchedulerExpression := "0 * * * *"

	lambdaRows := sqlmock.NewRows([]string{"id", "function_name", "function_arn"}).
		AddRow(1, util.GenLambdaFuncName(ctx, testFuncName), testFuncArn)
	mock.ExpectQuery(`SELECT \* FROM "lambda"`).
		WillReturnRows(lambdaRows)

	schedulerRows := sqlmock.NewRows([]string{"id", "lambda_id", "expression"}).
		AddRow(1, 1, testSchedulerExpression)
	mock.ExpectQuery(`SELECT \* FROM "lambda_scheduler"`).
		WillReturnRows(schedulerRows)

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	lambda, err := repo.FindByNameOrARN(ctx, testFuncName)

	assert.NoError(t, err)
	assert.Equal(t, util.GenLambdaFuncName(ctx, testFuncName), lambda.FunctionName)
	assert.Equal(t, testFuncArn, lambda.FunctionArn)
	assert.NotNil(t, lambda.Scheduler)
	assert.Equal(t, uint64(1), lambda.Scheduler.LambdaID)
	assert.Equal(t, testSchedulerExpression, lambda.Scheduler.Expression)
}

func TestFindByNameOrARNNotFound(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimIss.Str(), "test-org")
	ctx.Set(constant.ClaimSub.Str(), "test-account")
	testFuncName := "testFunc"

	lambdaRows := sqlmock.NewRows([]string{"id", "function_name", "function_arn"})
	mock.ExpectQuery(`SELECT \* FROM "lambda"`).
		WillReturnRows(lambdaRows)

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	lambda, err := repo.FindByNameOrARN(ctx, testFuncName)

	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("none lambda found by: %s", testFuncName), err.Error())
	assert.Nil(t, lambda)
}

func TestLambdaInfoSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimIss.Str(), "test-org")
	ctx.Set(constant.ClaimSub.Str(), "test-account")
	testFuncName := "testFunc"
	testFuncArn := "testArn"
	testSchedulerExpression := "0 * * * *"

	lambdaRows := sqlmock.NewRows([]string{"id", "function_name", "function_arn"}).
		AddRow(1, util.GenLambdaFuncName(ctx, testFuncName), testFuncArn)
	mock.ExpectQuery(`SELECT \* FROM "lambda"`).
		WillReturnRows(lambdaRows)

	schedulerRows := sqlmock.NewRows([]string{"id", "lambda_id", "expression"}).
		AddRow(1, 1, testSchedulerExpression)
	mock.ExpectQuery(`SELECT \* FROM "lambda_scheduler"`).
		WillReturnRows(schedulerRows)

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	lambdaInfo, err := repo.LambdaInfo(ctx, 1, testFuncName)

	assert.NoError(t, err)
	assert.Equal(t, util.GenLambdaFuncName(ctx, testFuncName), lambdaInfo.FunctionName)
	assert.Equal(t, testFuncArn, lambdaInfo.FunctionArn)
	assert.NotNil(t, lambdaInfo.Scheduler)
	assert.Equal(t, uint64(1), lambdaInfo.Scheduler.LambdaID)
	assert.Equal(t, testSchedulerExpression, lambdaInfo.Scheduler.Expression)
}

func TestLambdaInfoNotFound(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimIss.Str(), "test-org")
	ctx.Set(constant.ClaimSub.Str(), "test-account")
	testFuncName := "testFunc"

	lambdaRows := sqlmock.NewRows([]string{"id", "function_name", "function_arn"})
	mock.ExpectQuery(`SELECT \* FROM "lambda"`).
		WillReturnRows(lambdaRows)

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	lambdaInfo, err := repo.LambdaInfo(ctx, 1, testFuncName)

	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("none lambda found by: %s", testFuncName), err.Error())
	assert.Nil(t, lambdaInfo)
}

func TestFindByAccountSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimIss.Str(), "test-org")
	ctx.Set(constant.ClaimSub.Str(), "test-account")
	testFuncName := "testFunc"
	testFuncArn := "testArn"
	testSchedulerExpression := "0 * * * *"

	lambdaRows := sqlmock.NewRows([]string{"id", "function_name", "function_arn"}).
		AddRow(1, util.GenLambdaFuncName(ctx, testFuncName), testFuncArn)
	mock.ExpectQuery(`SELECT \* FROM "lambda"`).
		WillReturnRows(lambdaRows)

	schedulerRows := sqlmock.NewRows([]string{"id", "lambda_id", "expression"}).
		AddRow(1, 1, testSchedulerExpression)
	mock.ExpectQuery(`SELECT \* FROM "lambda_scheduler"`).
		WillReturnRows(schedulerRows)

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	lambdaInfos, err := repo.FindByAccount(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(lambdaInfos))
	assert.Equal(t, util.GenLambdaFuncName(ctx, testFuncName), lambdaInfos[0].FunctionName)
	assert.Equal(t, testFuncArn, lambdaInfos[0].FunctionArn)
	assert.NotNil(t, lambdaInfos[0].Scheduler)
	assert.Equal(t, uint64(1), lambdaInfos[0].Scheduler.LambdaID)
	assert.Equal(t, testSchedulerExpression, lambdaInfos[0].Scheduler.Expression)
}

func TestFindByAccountError(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimIss.Str(), "test-org")
	ctx.Set(constant.ClaimSub.Str(), "test-account")

	mock.ExpectQuery(`SELECT \* FROM "lambda"`).
		WillReturnError(errors.New("find record error"))

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	lambdaInfos, err := repo.FindByAccount(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("failed to query lambda by account, err: %s", "find record error"), err.Error())
	assert.Nil(t, lambdaInfos)
}

func TestPersistRegResultSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)

	mock.ExpectBegin()
	mock.ExpectCommit()

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	err := repo.PersistRegResult(ctx, func(tx *gorm.DB) error {
		return nil
	})

	assert.NoError(t, err)
}

func TestPersistRegResultError(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)

	mock.ExpectBegin()
	mock.ExpectRollback()

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	err := repo.PersistRegResult(ctx, func(tx *gorm.DB) error {
		return nil
	})

	assert.Error(t, err)
	assert.Equal(t, "failed to persist lambda registration result", err.Error())
}

func TestDeleteLambdaTXSuccess(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)

	mock.ExpectBegin()
	mock.ExpectCommit()

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	err := repo.DeleteLambdaTX(ctx, func(tx *gorm.DB) error {
		return nil
	})

	assert.NoError(t, err)
}

func TestDeleteLambdaTXError(t *testing.T) {
	sqldb, gormdb, mock := DbMock(t)
	defer sqldb.Close()

	ctx := new(gin.Context)

	mock.ExpectBegin()
	mock.ExpectRollback()

	repo := &lambda{
		Instance: &db.Instance{DB: gormdb},
	}
	err := repo.DeleteLambdaTX(ctx, func(tx *gorm.DB) error {
		return nil
	})

	assert.Error(t, err)
	assert.Equal(t, "failed to remove lambda", err.Error())
}
