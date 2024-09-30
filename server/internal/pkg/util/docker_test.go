package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFileSystem struct {
	mock.Mock
}

func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	return nil, args.Error(1)
}

func (m *mockFileSystem) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func TestIsRunningInsideDockerSuccess(t *testing.T) {
	mockFS := new(mockFileSystem)

	mockFS.On("Stat", "/"+".dockerenv").Return(nil, nil)

	assert.True(t, IsRunningInsideDocker(mockFS))
}

func TestIsRunningInsideDockerCGroupNotFound(t *testing.T) {
	mockFS := new(mockFileSystem)

	mockFS.On("Stat", "/.dockerenv").Return(nil, os.ErrNotExist)
	mockFS.On("ReadFile", "/proc/self/cgroup").Return(nil, os.ErrNotExist)

	assert.False(t, IsRunningInsideDocker(mockFS))
}

func TestIsRunningInsideDockerHaveCGroup(t *testing.T) {
	mockFS := new(mockFileSystem)

	mockFS.On("Stat", "/.dockerenv").Return(nil, os.ErrNotExist)

	mockCGroupContent := `
	11:name=systemd:/docker/3601745b3bd54d9780436faa5f0e4f72b98ae72d12e56c53c633f780f0d595b0
	`
	mockFS.On("ReadFile", "/proc/self/cgroup").Return([]byte(mockCGroupContent), nil)

	assert.True(t, IsRunningInsideDocker(mockFS))
}

func TestIsRunningInsideDockerHaveCGroupWithNoDocker(t *testing.T) {
	mockFS := new(mockFileSystem)

	mockFS.On("Stat", "/.dockerenv").Return(nil, os.ErrNotExist)
	mockFS.On("ReadFile", "/proc/self/cgroup").Return([]byte("invalid content"), nil)

	assert.False(t, IsRunningInsideDocker(mockFS))
}
