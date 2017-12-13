package main

import (
	"path/filepath"
	"strings"
	"os"
	"io"
	"archive/tar"
)

func Tar(source, target string) error {
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {

			if info.IsDir() {
				return nil
			}

			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			header.Name = strings.TrimPrefix(strings.TrimPrefix(path, source), "/")

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}
