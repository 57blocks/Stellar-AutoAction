package wallet

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/testdata"

	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestResourceCreateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("POST", "/create", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockWalletService(ctrl)

	mockResp := &dto.RespCreateWallet{
		Address: "test",
	}
	mockService.EXPECT().Create(ctx).Return(mockResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.Create(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)

	resp := &dto.RespCreateWallet{}
	err := json.Unmarshal(w.Body.Bytes(), resp)
	assert.Nil(t, err)
	assert.Equal(t, mockResp, resp)
}

func TestResourceCreateServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("POST", "/create", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockWalletService(ctrl)

	mockService.EXPECT().Create(ctx).Return(nil, errors.New("error"))

	cd := &resource{
		service: mockService,
	}

	cd.Create(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "error", ctx.Errors.Last().Error())
}

func TestResourceRemoveSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("DELETE", "/remove/test", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockWalletService(ctrl)

	mockService.EXPECT().Remove(ctx, gomock.Any()).Return(nil)

	cd := &resource{
		service: mockService,
	}

	cd.Remove(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)
}

func TestResourceRemoveServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("DELETE", "/remove/test", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockWalletService(ctrl)

	mockService.EXPECT().Remove(ctx, gomock.Any()).Return(errors.New("error"))

	cd := &resource{
		service: mockService,
	}

	cd.Remove(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "error", ctx.Errors.Last().Error())
}

func TestResourceListSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/list", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockWalletService(ctrl)

	mockResp := &dto.RespListWallets{
		Data: []dto.RespListWallet{
			{
				Address: "test",
			},
		},
	}
	mockService.EXPECT().List(ctx).Return(mockResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.List(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)

	resp := &dto.RespListWallets{}
	err := json.Unmarshal(w.Body.Bytes(), resp)
	assert.Nil(t, err)
	assert.Equal(t, mockResp, resp)
}

func TestResourceListServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/list", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockWalletService(ctrl)

	mockService.EXPECT().List(ctx).Return(nil, errors.New("error"))

	cd := &resource{
		service: mockService,
	}

	cd.List(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "error", ctx.Errors.Last().Error())
}

func TestResourceVerifySuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/verify/test", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockWalletService(ctrl)

	mockService.EXPECT().Verify(ctx, gomock.Any()).Return(nil, nil)

	cd := &resource{
		service: mockService,
	}

	cd.Verify(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)
}

func TestResourceVerifyServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/verify/test", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockWalletService(ctrl)

	mockService.EXPECT().Verify(ctx, gomock.Any()).Return(nil, errors.New("error"))

	cd := &resource{
		service: mockService,
	}

	cd.Verify(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "error", ctx.Errors.Last().Error())
}
