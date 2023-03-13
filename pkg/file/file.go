package file

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	timeNow := time.Now().Format("20060102150405")

	var tempFiles []string
	for key, value := range files {
		// path.Base(pathString)函数,pathString的值必须为linux风格的路径，即 "/" 才能够正常的获取最后的路径段的值。在如果路径是windows风格的，需要使用 pathfile.ToSlash()函数，将路径转为linux风格
		value = filepath.ToSlash(value)
		fileSuffix := path.Ext(value)
		newpath := path.Dir(value) + "/" + timeNow + "-" + fmt.Sprintf("%04d", key) + fileSuffix
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

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		// error handling
		return false
	}

	return fileInfo.IsDir()
}

func GetRedirectUrl(oldUrl string) (newUrl string) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(oldUrl)
	if err != nil {
		fmt.Println(err)
		newUrl = oldUrl
		return
	}

	newUrl = resp.Header.Get("location")
	if newUrl == "" {
		newUrl = oldUrl
		return
	}

	return
}

func Create(fileName string, fileContent any) (err error) {
	filePath := path.Dir(fileName)
	if !Exists(filePath) {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	outputFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer func(outputFile *os.File) {
		err = outputFile.Close()
		if err != nil {
			return
		}
	}(outputFile)
	encoder := json.NewEncoder(outputFile)
	err = encoder.Encode(fileContent)
	if err != nil {
		return
	}

	return
}
