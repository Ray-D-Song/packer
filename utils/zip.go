package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "dist/") {
			fpath := filepath.Join(dest, f.Name[len("dist/"):])
			if f.FileInfo().IsDir() {
				if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
					return err
				}
				continue
			}
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			w, err := os.Create(fpath)
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, rc); err != nil {
				return err
			}
			w.Close()
		}
	}
	return nil
}
