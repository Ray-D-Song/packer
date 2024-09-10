package handler

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"ray-d-song.com/packer/model"
	"ray-d-song.com/packer/utils"
)

func LibDownload(c *fiber.Ctx) error {
	dep := model.Dependency{}
	c.BodyParser(&dep)
	if !utils.CheckLibVersionExist(dep.Name, dep.Version) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Library version not found",
		})
	}
	depPath := path.Join(utils.StorageDir, dep.Name, dep.Version)
	zipData, err := zipFolder(depPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error zipping library",
		})
	}
	return c.Send(zipData)
}

func zipFolder(path string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		zipFile, err := zipWriter.Create(strings.TrimPrefix(filePath, path+"/"))
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
