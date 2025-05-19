package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/shreyansh-ML/movieapp/file/internal/storage"
)

type Local struct {
	basePath string
	maxSize  int64
}

func New(basePath string, maxSize int64) *Local {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		panic("failed to get absolute path: " + err.Error())
	}
	return &Local{
		basePath: absPath,
		maxSize:  maxSize,
	}
}
func (l *Local) Save(ctx context.Context, fileName string, data io.Reader) (string, error) {
	// Create the directory if it doesn't exist
	id, ok := ctx.Value("id").(int64)
	//id, ok := idValue.(int64)
	if !ok {
		return "", fmt.Errorf("invalid type for id, expected int64")
	}
	path := l.basePath + "/" + strconv.Itoa(id)
	err := os.MkdirAll(path, os.ModePerm) // Unix permission bits, 0o777

	if err != nil {
		return "", storage.ErrDirPermDenied
	}
	// Create the file
	filePath := path + "/" + fileName
	//file, err := os.Create(filePath)

	_, err = os.Stat(filePath)
	if err == nil {
		err = os.Remove(filePath)
		if err != nil {
			return "", fmt.Errorf("unable to remove old file %v", err)
		}
	} else if !os.IsNotExist(err) {
		// if this is anything other than a not exists error
		return "", fmt.Errorf("Unable to get file info: %w", err)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()

	// write the contents to the new file
	// ensure that we are not writing greater than max bytes
	_, err = io.Copy(f, data)
	if err != nil {
		return "", fmt.Errorf("Unable to write to file: %w", err)
	}

	return filePath, nil
}
