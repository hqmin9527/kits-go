package utils

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

func Exist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ReadObj(filePath string, obj any) error {
	data, err := ReadBytes(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}
	return nil
}

func WriteObj(filePath string, obj any) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}
	return WriteBytes(filePath, data)
}

func ReadBytes(filePath string) ([]byte, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}
	return b, nil
}

func WriteBytes(filePath string, data []byte) error {
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return errors.Wrap(err, "write file")
	}
	return nil
}

func CopyFile(src, dest string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return errors.Wrap(err, "stat src")
	}
	input, err := os.ReadFile(src)
	if err != nil {
		return errors.Wrap(err, "read src")
	}
	err = os.WriteFile(dest, input, fi.Mode())
	if err != nil {
		return errors.Wrap(err, "write dest")
	}
	return nil
}
