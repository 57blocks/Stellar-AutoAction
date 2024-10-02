package restyx

import (
	"net/http"
	"os"
	"testing"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// Before test, setup log
func TestMain(m *testing.M) {
	os.Setenv("CS_ENDPOINT", "https://api.fake.com")
	os.Setenv("CS_ORGANIZATION", "ORG1")
	config.Setup("../../config/")
	testConfig := config.Configuration{
		Log: config.Log{
			Level:    "debug",
			Encoding: "json",
		},
	}
	logx.Setup(&testConfig)

	os.Exit(m.Run())
}

func TestAddCSRoleSuccess(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	expectedRoleName := "test_role_name"
	expectedRoleId := "test_role_id"
	httpmock.RegisterResponder("POST", "https://api.fake.com/v0/org/ORG1/roles",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `{"name": "`+expectedRoleName+`", "role_id": "`+expectedRoleId+`"}`)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})
	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	resp, err := cd.AddCSRole(ctx, "test_cs_token", "test_org", "test_account")

	assert.NoError(t, err)
	assert.Equal(t, expectedRoleName, resp.Name)
	assert.Equal(t, expectedRoleId, resp.RoleId)
}

func TestAddCSRoleFailed(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.fake.com/v0/org/ORG1/roles",
		httpmock.NewStringResponder(400, `{"status":{"message": "error", "code": 400}}`))

	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	resp, err := cd.AddCSRole(ctx, "test_cs_token", "test_org", "test_account")

	assert.Error(t, err)
	assert.Equal(t, `create cube signer role occurred error: 400, {"status":{"message": "error", "code": 400}}`, err.Error())
	assert.Nil(t, resp)
}

func TestAddCSKeySuccess(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	expectedKeyID := "test_key_id"
	httpmock.RegisterResponder("POST", "https://api.fake.com/v0/org/ORG1/keys",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `{"keys":[{"key_id":"`+expectedKeyID+`"}]}`)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})
	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	key, err := cd.AddCSKey(ctx, "test_cs_token")

	assert.NoError(t, err)
	assert.Equal(t, expectedKeyID, key)
}

func TestAddCSKeyFailed(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.fake.com/v0/org/ORG1/keys",
		httpmock.NewStringResponder(400, `{"status":{"message": "error", "code": 400}}`))

	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	key, err := cd.AddCSKey(ctx, "test_cs_token")

	assert.Error(t, err)
	assert.Equal(t, `create cube signer key occurred error: 400, {"status":{"message": "error", "code": 400}}`, err.Error())
	assert.Equal(t, "", key)
}

func TestAddCSKeyToRoleSuccess(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	expectedKeyID := "test_key_id"
	httpmock.RegisterResponder("PUT", "https://api.fake.com/v0/org/ORG1/roles/Role1/add_keys",
		httpmock.NewStringResponder(200, `{"status":{"message": "success", "code": 200}}`))
	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	err := cd.AddCSKeyToRole(ctx, "test_cs_token", expectedKeyID, "Role1")

	assert.NoError(t, err)
}

func TestAddCSKeyToRoleFailed(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", "https://api.fake.com/v0/org/ORG1/roles/Role1/add_keys",
		httpmock.NewStringResponder(400, `{"status":{"message": "error", "code": 400}}`))
	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	err := cd.AddCSKeyToRole(ctx, "test_cs_token", "test_key_id", "Role1")

	assert.Error(t, err)
	assert.Equal(t, `add cube signer key to role occurred error: 400, {"status":{"message": "error", "code": 400}}`, err.Error())
}

func TestDeleteCSKeySuccess(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.fake.com/v0/org/ORG1/keys/Key1",
		httpmock.NewStringResponder(200, `{"status":{"message": "success", "code": 200}}`))
	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	err := cd.DeleteCSKey(ctx, "test_cs_token", "Key1")

	assert.NoError(t, err)
}

func TestDeleteCSKeyFailed(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.fake.com/v0/org/ORG1/keys/Key1",
		httpmock.NewStringResponder(400, `{"status":{"message": "error", "code": 400}}`))
	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	err := cd.DeleteCSKey(ctx, "test_cs_token", "Key1")

	assert.Error(t, err)
	assert.Equal(t, `delete cube signer key occurred error: 400, {"status":{"message": "error", "code": 400}}`, err.Error())
}

func TestDeleteCSKeyFromRoleSuccess(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.fake.com/v0/org/ORG1/roles/Role1/keys/Key1",
		httpmock.NewStringResponder(200, `{"status":{"message": "success", "code": 200}}`))
	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	err := cd.DeleteCSKeyFromRole(ctx, "test_cs_token", "Key1", "Role1")

	assert.NoError(t, err)
}

func TestDeleteCSKeyFromRoleFailed(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.fake.com/v0/org/ORG1/roles/Role1/keys/Key1",
		httpmock.NewStringResponder(400, `{"status":{"message": "error", "code": 400}}`))
	ctx := new(gin.Context)

	cd := &restyx{client: restyClient}
	err := cd.DeleteCSKeyFromRole(ctx, "test_cs_token", "Key1", "Role1")

	assert.Error(t, err)
	assert.Equal(t, `delete cube signer key from role occurred error: 400, {"status":{"message": "error", "code": 400}}`, err.Error())
}
