package file

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func GetNameWithoutExt() string {
	_, fullFilename, _, _ := runtime.Caller(1)
	//获取文件名带后缀
	filenameWithSuffix := path.Base(fullFilename)
	//获取文件后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	//获取文件名
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)

	return filenameOnly
}

func GetFiles(pathName string) (files []string) {
	err := filepath.Walk(pathName, func(pathFile string, info os.FileInfo, err error) error {
		if pathName == pathFile {
			return nil
		}

		files = append(files, pathFile)
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func SerialRename(files []string) {
	new := time.Now().Format("20060102150405")

	var tempFiles []string
	for key, value := range files {
		fileSuffix := path.Ext(value)
		newpath := path.Dir(value) + "/" + new + "-" + fmt.Sprintf("%04d", key) + fileSuffix
		err := os.Rename(value, newpath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		tempFiles = append(tempFiles, newpath)
	}

	for key, value := range tempFiles {
		fileSuffix := path.Ext(value)
		newpath := path.Dir(value) + "/" + fmt.Sprintf("%04d", key) + fileSuffix
		err := os.Rename(value, newpath)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return
}

func Exists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}
