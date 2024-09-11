package eslint

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/util"

	"github.com/google/uuid"
)

const (
	LintDir    = "internal/third-party/eslint"
	LintConfig = LintDir + "/eslint.config.mjs"
	JSExt      = ".js"
)

func init() {
	if util.IsRunningInsideDocker() {
		output, err := exec.Command("npx", "eslint", "-c", LintConfig, ".").CombinedOutput()
		if err != nil {
			panic(fmt.Sprintf("failed to init eslint: %s\n%s", string(output), err.Error()))
		}
	}
}

func Check(zipFile *multipart.FileHeader) error {
	tmpFilePath := filepath.Join(".", LintDir, uuid.New().String(), zipFile.Filename)
	tmpDir := filepath.Dir(tmpFilePath)
	defer os.RemoveAll(tmpDir)

	if err := os.MkdirAll(filepath.Dir(tmpFilePath), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	open, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer open.Close()

	if _, err := io.Copy(out, open); err != nil {
		return err
	}

	if err := exec.Command("unzip", tmpFilePath, "-d", tmpDir).Run(); err != nil {
		return err
	}

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), JSExt) {
			output, err := exec.Command("npx", "eslint", "-c", LintConfig, filepath.Join(tmpDir, file.Name())).
				CombinedOutput()
			if err != nil {
				errMessage := strings.Trim(string(output), "\n")
				if strings.Contains(errMessage, "/") {
					fileName := strings.Split(errMessage, "/")[len(strings.Split(errMessage, "/"))-1]

					return errorx.BadRequest(fmt.Sprintf("%s/%s\n error: %s", zipFile.Filename, fileName, errMessage))
				} else {
					return errorx.BadRequest(errMessage)
				}
			}
		}
	}

	return nil
}
