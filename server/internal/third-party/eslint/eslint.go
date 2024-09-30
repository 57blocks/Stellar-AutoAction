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
	LintDir       = "internal/third-party/eslint"
	LintConfig    = LintDir + "/eslint.config.mjs"
	JSExt         = ".js"
	PathSeparator = string(os.PathSeparator)
)

var fileSystem util.FileSystem = &util.RealFileSystem{}

func init() {
	if util.IsRunningInsideDocker(fileSystem) {
		output, err := exec.Command("npx", "eslint", "-c", LintConfig, ".").CombinedOutput()
		if err != nil {
			panic(fmt.Sprintf("failed to init eslint: %s%s%s", string(output), PathSeparator, err.Error()))
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
	lineBreak := newLine()
	for _, file := range files {
		if strings.HasSuffix(file.Name(), JSExt) {
			output, err := exec.Command("npx", "eslint", "-c", LintConfig, filepath.Join(tmpDir, file.Name())).
				CombinedOutput()
			errMessage := strings.TrimSpace(string(output))
			if err == nil {
				continue
			}
			if !strings.Contains(errMessage, lineBreak) {
				return errorx.BadRequest(errMessage)
			}
			messageLines := strings.Split(errMessage, lineBreak)
			for _, line := range messageLines {
				if strings.Contains(line, ":") && strings.Contains(line, "error") {
					errMessage = strings.TrimSpace(line)
					break
				}
			}
			return errorx.BadRequest(fmt.Sprintf("%s/%s: %s", zipFile.Filename, file.Name(), errMessage))
		}
	}

	return nil
}

func newLine() string {
	if os.PathSeparator == '/' {
		return "\n"
	}
	return "\r\n"
}
