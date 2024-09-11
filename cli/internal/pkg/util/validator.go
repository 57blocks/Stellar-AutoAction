package util

import (
	"archive/zip"
	"os"
	"strings"

	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
)

func ValidateZipFiles(paths []string) error {
	for _, path := range paths {
		if err := CheckJsLambdaZipFile(path); err != nil {
			return err
		}
	}
	return nil
}

func CheckJsLambdaZipFile(path string) error {
	zipFile, err := os.Open(path)
	// check if the file exists
	if err != nil {
		return errorx.BadRequest(err.Error())
	}
	defer zipFile.Close()

	stat, err := zipFile.Stat()
	if err != nil {
		return errorx.Internal(err.Error())
	}

	// check if the file is a valid zip file
	archive, err := zip.NewReader(zipFile, stat.Size())
	if err != nil {
		return errorx.Internal(err.Error() + ": " + path)
	}

	// check if the zip file contains any js files, will use ESLint to check the content at the sever side
	validJsFiles := make([]string, 0)
	for _, file := range archive.File {
		if strings.Contains(file.Name, "/") {
			continue
		}
		if strings.HasSuffix(file.Name, ".js") {
			validJsFiles = append(validJsFiles, file.Name)
		}
	}
	if len(validJsFiles) == 0 {
		return errorx.BadRequest("no valid js files found in: " + path)
	}

	return nil
}
