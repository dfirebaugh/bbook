package web

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed static
var StaticFiles embed.FS

func CopyStaticFiles(srcFS embed.FS, srcDir, dstDir string) error {
	err := fs.WalkDir(srcFS, srcDir, func(srcPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(srcDir, srcPath)
		dstPath := filepath.Join(dstDir, relPath)

		err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
		if err != nil {
			return err
		}

		if filepath.Base(srcPath) == "SUMMARY.md" {
			return nil
		}

		if filepath.Ext(srcPath) == ".md" {
			return nil
		}

		srcFile, err := srcFS.Open(srcPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		return copyFile(srcFile, dstPath)
	})

	return err
}

func copyFile(srcFile fs.File, dstPath string) error {
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}
