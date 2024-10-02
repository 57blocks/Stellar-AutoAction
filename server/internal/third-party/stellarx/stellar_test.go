package stellarx

import (
	"errors"
	"testing"

	"github.com/57blocks/auto-action/server/internal/testdata"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
)

func TestAccountDetailSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := testdata.NewMockHorizonClient(ctrl)

	ctx := new(gin.Context)
	expectedAccount := horizon.Account{
		AccountID: "test_account_id",
		Sequence:  123456,
	}
	req := horizonclient.AccountRequest{AccountID: expectedAccount.AccountID}

	mockClient.EXPECT().AccountDetail(req).Return(expectedAccount, nil)

	s := &stellar{
		client: mockClient,
	}

	account, err := s.AccountDetail(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedAccount, account)
}

func TestAccountDetailError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := testdata.NewMockHorizonClient(ctrl)

	ctx := new(gin.Context)
	req := horizonclient.AccountRequest{AccountID: "test_account_id"}

	mockClient.EXPECT().AccountDetail(req).Return(horizon.Account{}, errors.New("error"))

	s := &stellar{
		client: mockClient,
	}

	account, err := s.AccountDetail(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "error", err.Error())
	assert.Equal(t, horizon.Account{}, account)
}
