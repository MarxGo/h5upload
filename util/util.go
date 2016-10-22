package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func CompleteDirPath(dirPath string) {
	os.MkdirAll(dirPath, 777)
}

func CheckFileMd5(filePath, md5Value string) (bool, error) {
	// if server has this file ,skip it
	// check current file's md5 is equal to request
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return false, err
	}

	md5Ctx := md5.New()
	io.Copy(md5Ctx, file)
	fileMd5 := fmt.Sprintf("%x", md5Ctx.Sum([]byte("")))
	return fileMd5 == md5Value, nil
}
