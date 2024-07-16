package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip extracts a zip archive to a specified destination.
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// MovePDF moves PDF files from the source to the destination.
func MovePDF(srcDir, destDir string) error {
	err := os.MkdirAll(filepath.Dir(destDir), os.ModePerm)
	if err != nil {
		return err
	}

	pdfFiles, err := filepath.Glob(filepath.Join(srcDir, "*.pdf"))
	if err != nil {
		return err
	}

	for _, pdfFile := range pdfFiles {
		destPath := filepath.Join(destDir, filepath.Base(pdfFile))
		err := os.Rename(pdfFile, destPath)
		if err != nil {
			return err
		}
	}

	return nil
}
