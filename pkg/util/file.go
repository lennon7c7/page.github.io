package util

import (
	"archive/zip"
	"io"
	"os"
	"page.github.io/pkg/config"
	"path/filepath"
	"strings"
)

// PathExists 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetCurrentDirectory 获取当前执行文件所在的目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {

	}
	return strings.Replace(dir, "\\", "/", -1)
}

// Getwd 获取当前执行文件所在的目录，如果是go run方式运行，则获取到执行go run命令所在的目录
func Getwd() string {
	dir, _ := os.Getwd()
	return dir
}

// 某目录打包成zip文件
// `./docs`, `oss.zip`
// srcDir 源目录
// zipFileName 打包后的路径
// downloadURL 下载链接
func ArchiveFolderToZip(srcDir string, zipFileName string) (downloadURL string, err error) {
	// 因为要预防旧文件无法覆盖，所以索性直接删除
	_ = os.RemoveAll(zipFileName)

	// 创建：zip文件
	zipFile, _ := os.Create(zipFileName)
	defer func() {
		_ = zipFile.Close()
	}()

	// 打开：zip文件
	archive := zip.NewWriter(zipFile)
	defer func() {
		_ = archive.Close()
	}()

	// 遍历路径信息
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, _ error) error {
		// 如果是源路径，提前进行下一个遍历
		if path == srcDir {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, srcDir+`\`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer func() {
				_ = file.Close()
			}()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	newPath := "static/" + zipFileName
	err = os.Rename(zipFileName, newPath)
	if err != nil {
		return
	}

	downloadURL = config.CFG.Domain.Api + "/static/" + zipFileName

	// 直接删除源目录
	_ = os.RemoveAll(srcDir)

	return
}

// @tarFile：压缩文件路径
// @dest：解压文件夹
func DeCompressByPath(tarFile string) (descFolder string, err error) {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return
	}

	defer func() {
		_ = srcFile.Close()
	}()

	return DeCompress(srcFile)
}

// @zipFile：压缩文件
// @dest：解压之后文件保存路径
func DeCompress(srcFile *os.File) (descFolder string, err error) {
	zipFile, err := zip.OpenReader(srcFile.Name())
	if err != nil {
		return
	}

	defer func() {
		_ = zipFile.Close()
	}()

	for index, innerFile := range zipFile.File {
		if index == 0 && innerFile.Name != "" {
			descFolder = innerFile.Name
		}

		fileInfo := innerFile.FileInfo()
		if fileInfo.IsDir() {
			err = os.MkdirAll(innerFile.Name, os.ModePerm)
			if err != nil {
				return
			}

			continue
		}

		readCloser, err := innerFile.Open()
		if err != nil {
			continue
		}

		newFile, err := os.Create(innerFile.Name)
		if err != nil {
			continue
		}

		_, err = io.Copy(newFile, readCloser)
		if err != nil {
			return descFolder, err
		}

		_ = readCloser.Close()
		_ = newFile.Close()
	}

	return
}
