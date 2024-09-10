package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"ray-d-song.com/packer/model"
)

type Response struct {
	Code    int16       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseOK(data any) fiber.Map {
	return fiber.Map{
		"code": 200,
		"data": data,
	}
}

func ResponseErr(code int16, message string) fiber.Map {
	return fiber.Map{
		"code":    code,
		"message": message,
	}
}

func Download(registry string, dep string) {
	parts := strings.Split(dep, "@")
	if len(parts) != 2 {
		return
	}
	dependency := model.Dependency{
		Name:    parts[0],
		Version: parts[1],
	}

	payload, err := json.Marshal(dependency)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("%s/api/lib/download", registry)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Error downloading dependency: %s", response.Status))
	}

	depPath := path.Join(GetLibsPath(), dep)
	err = EnsureDirExists(depPath)
	if err != nil {
		panic(err)
	}

	dst, err := os.CreateTemp("", "uploaded-*.zip")
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(dst, response.Body); err != nil {
		panic(err)
	}
	defer os.Remove(dst.Name())
	if err := Unzip(dst.Name(), depPath); err != nil {
		panic(err)
	}
}
