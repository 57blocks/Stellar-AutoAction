package lambda

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/testdata"

	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestResourceRegisterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))

	writer.WriteField("expression", "test expression")
	writer.WriteField("payload", "test payload")
	writer.Close()

	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockRegistResp := []*dto.RespRegister{
		{
			Lambda: &dto.RespLamBrief{
				Name: "test-func",
			},
			Scheduler: &dto.RespSchBrief{
				Name: "test-scheduler",
			},
		},
	}
	mockService.EXPECT().Register(ctx, gomock.Any()).Return(mockRegistResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.Register(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, ctx.Errors)

	var actualResp []*dto.RespRegister
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, mockRegistResp[0].Lambda.Name, actualResp[0].Lambda.Name)
	assert.Equal(t, mockRegistResp[0].Scheduler.Name, actualResp[0].Scheduler.Name)
}

func TestResourceRegisterParseMultipartFormError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("POST", "/register", nil)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req

	cd := &resource{}

	cd.Register(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "failed to parse multipart form: request Content-Type isn't multipart/form-data", ctx.Errors.Last().Error())
}

func TestResourceRegisterServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))

	writer.WriteField("expression", "test expression")
	writer.WriteField("payload", "test payload")
	writer.Close()

	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockService.EXPECT().Register(ctx, gomock.Any()).
		Return(nil, errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.Register(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}

func TestResourceInvokeSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jsonBody := []byte(`{ "key": "value" }`)
	req := httptest.NewRequest("POST", "/invoke/test-func", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockInvokeResp := &dto.RespInvoke{
		StatusCode: 204,
	}
	mockService.EXPECT().Invoke(ctx, gomock.Any()).Return(mockInvokeResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.Invoke(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)

	var actualResp *dto.RespInvoke
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)

	assert.Equal(t, mockInvokeResp.StatusCode, actualResp.StatusCode)
}

func TestResourceInvokeBindJSONError(t *testing.T) {
	req := httptest.NewRequest("POST", "/invoke/test-func", nil)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req

	cd := &resource{}

	cd.Invoke(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "EOF", ctx.Errors.Last().Error())
}

func TestResourceInvokeServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jsonBody := []byte(`{ "key": "value" }`)
	req := httptest.NewRequest("POST", "/invoke/test-func", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockService.EXPECT().Invoke(ctx, gomock.Any()).Return(nil, errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.Invoke(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}

func TestResourceListSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/list?full=true", nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockListResp := []*dto.RespInList{
		{
			FunctionName: "test-func",
		},
	}

	mockService.EXPECT().List(ctx, gomock.Any()).Return(mockListResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.List(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)

	var actualResp []*dto.RespInList
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, mockListResp[0].FunctionName, actualResp[0].FunctionName)
}

func TestResourceListSBindQueryError(t *testing.T) {
	req := httptest.NewRequest("GET", "/list?full=invalid", nil)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req

	cd := &resource{}

	cd.List(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, `strconv.ParseBool: parsing "invalid": invalid syntax`, ctx.Errors.Last().Error())
}

func TestResourceListServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/list?full=true", nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockService.EXPECT().List(ctx, gomock.Any()).Return(nil, errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.List(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}

func TestResourceInfoSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/info/test-func", nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockInfoResp := &dto.RespInfo{
		FunctionName: "test-func",
	}

	mockService.EXPECT().Info(ctx, gomock.Any()).Return(mockInfoResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.Info(ctx)
	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	var actualResp *dto.RespInfo
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, mockInfoResp.FunctionName, actualResp.FunctionName)
}

func TestResourceInfoServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/info/test-func", nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockService.EXPECT().Info(ctx, gomock.Any()).Return(nil, errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.Info(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}

func TestResourceLogsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/logs/test-func", nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockService.EXPECT().Logs(ctx, gomock.Any(), gomock.Any()).Return(nil)

	cd := &resource{
		service: mockService,
	}

	cd.Logs(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)
}

func TestResourceLogsServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("GET", "/logs/test-func", nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockService.EXPECT().Logs(ctx, gomock.Any(), gomock.Any()).Return(errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.Logs(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}

func TestResourceRemoveSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("DELETE", "/remove/test-func", nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockRemoveResp := &dto.RespRemove{
		Lambdas: dto.RespLamBrief{
			Name: "test-func",
		},
	}
	mockService.EXPECT().Remove(ctx, gomock.Any()).Return(mockRemoveResp, nil)

	cd := &resource{
		service: mockService,
	}

	cd.Remove(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, ctx.Errors)

	var actualResp *dto.RespRemove
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, mockRemoveResp.Lambdas.Name, actualResp.Lambdas.Name)
}

func TestResourceRemoveServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest("DELETE", "/remove/test-func", nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService := testdata.NewMockLambdaService(ctrl)

	mockService.EXPECT().Remove(ctx, gomock.Any()).Return(nil, errors.New("service error"))

	cd := &resource{
		service: mockService,
	}

	cd.Remove(ctx)

	assert.NotNil(t, ctx.Errors)
	assert.Equal(t, "service error", ctx.Errors.Last().Error())
}
