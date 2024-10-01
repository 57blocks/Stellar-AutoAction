package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestResourceSignupSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jsonBody := []byte(`{ "account": "foo", "organization": "bar", "password": "baz" }`)
	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := NewMockService(ctrl)

	mockService.EXPECT().Signup(ctx, gomock.Any()).Return(nil)

	cd := &resource{
		service: mockService,
	}

	cd.Signup(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)
}

func TestResourceSignupBindJSONError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("POST", "/signup", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	cd := &resource{
		service: NewMockService(ctrl),
	}

	cd.Signup(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "EOF", ctx.Errors.Last().Error())
}

func TestResourceSignupServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jsonBody := []byte(`{ "account": "foo", "organization": "bar", "password": "baz" }`)
	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := NewMockService(ctrl)
	mockService.EXPECT().Signup(ctx, gomock.Any()).Return(errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.Signup(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}

func TestResourceLoginSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jsonBody := []byte(`{ "account": "foo", "password": "baz" }`)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := NewMockService(ctrl)

	expectedResp := &dto.RespCredential{
		Account:      "test-account",
		Organization: "test-org",
	}
	mockService.EXPECT().Login(ctx, gomock.Any()).Return(expectedResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.Login(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)

	resp := &dto.RespCredential{}
	err := json.Unmarshal(w.Body.Bytes(), resp)
	assert.Nil(t, err)
	assert.Equal(t, expectedResp, resp)
}

func TestResourceLoginBindJSONError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	cd := &resource{
		service: NewMockService(ctrl),
	}

	cd.Login(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "EOF", ctx.Errors.Last().Error())
}

func TestResourceLoginServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jsonBody := []byte(`{ "account": "foo", "password": "baz" }`)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := NewMockService(ctrl)
	mockService.EXPECT().Login(ctx, gomock.Any()).Return(nil, errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.Login(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}

func TestResourceLogoutSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("DELETE", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constant.ClaimRaw.Str(), "test-raw")
	ctx.Request = req

	mockService := NewMockService(ctrl)

	expectedResp := &dto.RespLogout{}
	mockService.EXPECT().Logout(ctx, gomock.Any()).Return(expectedResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.Logout(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)

	resp := &dto.RespLogout{}
	err := json.Unmarshal(w.Body.Bytes(), resp)
	assert.Nil(t, err)
	assert.Equal(t, expectedResp, resp)
}

func TestResourceLogoutUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("DELETE", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	cd := &resource{
		service: NewMockService(ctrl),
	}

	cd.Logout(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "request unauthorized", ctx.Errors.Last().Error())
}

func TestResourceLogoutRawNotString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("DELETE", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constant.ClaimRaw.Str(), 123)
	ctx.Request = req

	cd := &resource{
		service: NewMockService(ctrl),
	}

	cd.Logout(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "raw is not string", ctx.Errors.Last().Error())
}

func TestResourceLogoutServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("DELETE", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constant.ClaimRaw.Str(), "test-raw")
	ctx.Request = req

	mockService := NewMockService(ctrl)
	mockService.EXPECT().Logout(ctx, gomock.Any()).Return(nil, errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.Logout(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}

func TestResourceRefreshSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constant.ClaimRaw.Str(), "test-raw")
	ctx.Request = req

	mockService := NewMockService(ctrl)

	expectedResp := &dto.RespCredential{}
	mockService.EXPECT().Refresh(ctx, gomock.Any()).Return(expectedResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.Refresh(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)

	resp := &dto.RespCredential{}
	err := json.Unmarshal(w.Body.Bytes(), resp)
	assert.Nil(t, err)
	assert.Equal(t, expectedResp, resp)
}

func TestResourceRefreshUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	cd := &resource{
		service: NewMockService(ctrl),
	}

	cd.Refresh(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "request unauthorized", ctx.Errors.Last().Error())
}

func TestResourceRefreshRawNotString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constant.ClaimRaw.Str(), 123)
	ctx.Request = req

	cd := &resource{
		service: NewMockService(ctrl),
	}

	cd.Refresh(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "raw is not string", ctx.Errors.Last().Error())
}

func TestResourceRefreshServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constant.ClaimRaw.Str(), "test-raw")
	ctx.Request = req

	mockService := NewMockService(ctrl)
	mockService.EXPECT().Refresh(ctx, gomock.Any()).Return(nil, errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.Refresh(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}
