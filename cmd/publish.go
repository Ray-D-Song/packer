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

	"github.com/spf13/viper"
	"ray-d-song.com/packer/dict"
)

func Publish() {
	name := viper.GetString("name")
	version := viper.GetString("version")

	permFilePath := filepath.Join(dict.PackerDir, "perm.key")
	permData, err := os.ReadFile(permFilePath)
	if err != nil {
		log.Fatalf("Failed to read perm key file: %v", err)
		return
	}
	perm := string(permData)

	ignores := viper.GetStringSlice("ignore")
	// Get all files in the current directory and subdirectories
	var files []string
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
		return
	}

	// Filter out ignored files
	var filteredFiles []string
	for _, file := range files {
		ignore := false
		for _, pattern := range ignores {
			matched, err := filepath.Match(pattern, file)
			if err != nil {
				log.Fatalf("Failed to match pattern: %v", err)
				return
			}
			if matched {
				ignore = true
				break
			}
		}
		if !ignore {
			filteredFiles = append(filteredFiles, file)
		}
	}

	// Zip the remaining files
	zipData, err := zipFiles(filteredFiles)
	if err != nil {
		log.Fatalf("Failed to zip folder: %v", err)
		return
	}

	// 获取压缩包大小
	size := len(zipData)

	// 预请求
	preCheckPayload := map[string]interface{}{
		"size":    size,
		"name":    name,
		"version": version,
		"perm":    perm,
	}
	preCheckBody, err := json.Marshal(preCheckPayload)
	if err != nil {
		log.Fatalf("Failed to marshal pre-check payload: %v", err)
		return
	}

	registry := viper.GetString("registry")
	preCheckURL := fmt.Sprintf("%s/api/lib/check", registry)
	preCheckResp, err := http.Post(preCheckURL, "application/json", bytes.NewBuffer(preCheckBody))
	if err != nil {
		log.Fatalf("Failed to perform pre-check: %v", err)
		return
	}
	defer preCheckResp.Body.Close()

	if preCheckResp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed: %v", preCheckResp.Status)
		return
	}

	var preCheckResult map[string]interface{}
	if err := json.NewDecoder(preCheckResp.Body).Decode(&preCheckResult); err != nil {
		log.Fatalf("Failed to decode pre-check response: %v", err)
		return
	}
	if code, ok := preCheckResult["code"].(int); ok && code != 200 {
		message, _ := preCheckResult["message"].(string)
		log.Fatalf("Pre-check failed: %s", message)
		return
	}
	bodyField, ok := preCheckResult["data"]
	if !ok || bodyField == nil {
		log.Fatalf("Invalid pre-check response: missing body")
		return
	}

	bodyMap, ok := bodyField.(map[string]interface{})
	if !ok {
		log.Fatalf("Invalid pre-check response: body is not a map")
		return
	}

	ticket, ok := bodyMap["ticket"].(string)
	if !ok {
		log.Fatalf("Invalid pre-check response: missing ticket")
		return
	}

	// 创建表单数据
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "upload.zip")
	if err != nil {
		log.Fatalf("Failed to create form file: %v", err)
		return
	}
	part.Write(zipData)
	writer.WriteField("name", name)
	writer.WriteField("version", version)
	writer.WriteField("size", fmt.Sprintf("%d", size))
	writer.Close()

	// 发送发布请求
	publishURL := fmt.Sprintf("%s/api/lib/publish?size=%d&name=%s&version=%s&ticket=%s", registry, size, name, version, ticket)
	req, err := http.NewRequest("POST", publishURL, body)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to publish: %v", resp.Status)
		return
	}

	log.Println("Publish successful")
}

func zipFiles(files []string) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		w, err := zipWriter.Create(file)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(w, f); err != nil {
			return nil, err
		}
	}
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
