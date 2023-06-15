package util

import (
	"archive/zip"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"os"
	"page.github.io/pkg/config"
	"path/filepath"
	"strings"
)

// 上传文件
func OssPutObject(name string, reader io.Reader) (err error) {
	// 创建OSSClient实例。
	client, err := oss.New(config.CFG.Aliyun.Oss.Endpoint, config.CFG.Aliyun.Oss.AccessKeyID, config.CFG.Aliyun.Oss.AccessKeySecret)
	if err != nil {
		log.Println("Oss Error", err)
		return err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(config.CFG.Aliyun.Oss.BucketName)
	if err != nil {
		log.Println("Oss Error1", err)
		return err
	}

	// 上传文件流。
	err = bucket.PutObject(name, reader)
	if err != nil {
		log.Println("Oss Error2", err)
		return err
	}
	return
}

// 上传本地文件
func OssPutObjectFromFile(name string, filePath string) (signedURL string, err error) {
	// 创建OSSClient实例。
	client, err := oss.New(config.CFG.Aliyun.Oss.Endpoint, config.CFG.Aliyun.Oss.AccessKeyID, config.CFG.Aliyun.Oss.AccessKeySecret)
	if err != nil {
		return
	}

	// 获取存储空间。
	bucket, err := client.Bucket(config.CFG.Aliyun.Oss.BucketName)
	if err != nil {
		return
	}

	// 上传文件流。
	err = bucket.PutObjectFromFile(name, filePath)
	if err != nil {
		return
	}

	signedURL, err = OssGetSignURL(name)

	return
}

// 获取签名后的链接
func OssGetSignURL(unsignedURL string) (signedURL string, err error) {
	// 创建OSSClient实例。
	client, err := oss.New(config.CFG.Aliyun.Oss.Endpoint, config.CFG.Aliyun.Oss.AccessKeyID, config.CFG.Aliyun.Oss.AccessKeySecret)
	if err != nil {
		log.Println("Oss New Error: ", err)
		return
	}

	// 获取存储空间
	bucket, err := client.Bucket(config.CFG.Aliyun.Oss.BucketName)
	if err != nil {
		log.Println("Oss Bucket Error: ", err)
		return
	}

	// 签名直传
	signedURL, err = bucket.SignURL(unsignedURL, oss.HTTPGet, 3600)
	if err != nil {
		log.Println("Oss SignURL Error: ", err)
		return
	}

	return
}

// 下载到本地文件
func OssGetObjectToFile(unsignedURL string, filePath string) (err error) {
	// 创建OSSClient实例。
	client, err := oss.New(config.CFG.Aliyun.Oss.Endpoint, config.CFG.Aliyun.Oss.AccessKeyID, config.CFG.Aliyun.Oss.AccessKeySecret)
	if err != nil {
		log.Println("Oss New Error: ", err)
		return
	}

	// 获取存储空间
	bucket, err := client.Bucket(config.CFG.Aliyun.Oss.BucketName)
	if err != nil {
		log.Println("Oss Bucket Error: ", err)
		return
	}

	err = bucket.GetObjectToFile(unsignedURL, filePath)
	if err != nil {
		log.Println("Oss GetObjectToFile Error: ", err)
		return
	}

	return
}

// 打包成zip文件
// `./docs`, `oss.zip`
// srcDir 源目录
// zipFileName 打包后的路径
func ArchiveFileToZip(srcDir string, zipFileName string) (err error) {
	// 预防：旧文件无法覆盖
	_ = os.RemoveAll(zipFileName)

	// 创建：zip文件
	zipfile, _ := os.Create(zipFileName)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

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
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}

		return nil
	})

	return
}
