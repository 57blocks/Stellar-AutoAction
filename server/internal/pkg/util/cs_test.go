package util

import (
	"os"
	"testing"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/stretchr/testify/assert"
)

// Before test, setup log
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

func TestGetAddressFromCSKeySuccess(t *testing.T) {
	address, err := GetAddressFromCSKey("Key#Stellar_test-key")
	assert.NoError(t, err)
	assert.Equal(t, "test-key", address)
}

func TestGetAddressFromCSKeyFailed(t *testing.T) {
	address, err := GetAddressFromCSKey("test-key")
	assert.Error(t, err)
	assert.Equal(t, "invalid key", err.Error())
	assert.Equal(t, "", address)
}

func TestGetCSKeyFromAddressSuccess(t *testing.T) {
	assert.Equal(t, "Key#Stellar_test-key", GetCSKeyFromAddress("test-key"))
}
