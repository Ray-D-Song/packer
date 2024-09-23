package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Publish() {
	name := viper.GetString("name")
	version := viper.GetString("version")

	// 压缩文件夹
	zipData, err := zipFolder(".")
	if err != nil {
		log.Fatalf("Failed to zip folder: %v", err)
	}

	// 获取压缩包大小
	size := len(zipData)

	// 预请求
	preCheckPayload := map[string]interface{}{
		"size":    size,
		"name":    name,
		"version": version,
	}
	preCheckBody, err := json.Marshal(preCheckPayload)
	if err != nil {
		log.Fatalf("Failed to marshal pre-check payload: %v", err)
	}

	registry := viper.GetString("registry")
	preCheckURL := fmt.Sprintf("%s/api/lib/check", registry)
	preCheckResp, err := http.Post(preCheckURL, "application/json", bytes.NewBuffer(preCheckBody))
	if err != nil {
		log.Fatalf("Failed to perform pre-check: %v", err)
	}
	defer preCheckResp.Body.Close()

	if preCheckResp.StatusCode != http.StatusOK {
		log.Fatalf("Pre-check failed: %v", preCheckResp.Status)
	}

	// 创建表单数据
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "upload.zip")
	if err != nil {
		log.Fatalf("Failed to create form file: %v", err)
	}
	part.Write(zipData)
	writer.WriteField("name", name)
	writer.WriteField("version", version)
	writer.WriteField("size", fmt.Sprintf("%d", size))
	writer.Close()

	// 发送发布请求
	publishURL := fmt.Sprintf("%s/api/lib/publish?size=%d&name=%s&version=%s", registry, size, name, version)
	req, err := http.NewRequest("POST", publishURL, body)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to publish: %v", resp.Status)
	}

	log.Println("Publish successful")
}

func zipFolder(folderPath string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := filepath.Walk(folderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		zipFile, err := zipWriter.Create(strings.TrimPrefix(filePath, folderPath+"/"))
		if err != nil {
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipFile, file)
		return err
	})

	if err != nil {
		return nil, err
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
