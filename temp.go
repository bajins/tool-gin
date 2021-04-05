package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var TempDirPath string

func destroyTemp(path string) error {
	if strings.Contains(path, "temp") {

	}
	//err := filepath.WalkDir(path, func(path string, fi os.DirEntry, err error) error {
	err := filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if !fi.IsDir() {
			err := os.Remove(path)
			if err != nil {
				return err
			}
			return nil
		}
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func DestroyTempDir() {
	err := os.RemoveAll(TempDirPath)
	if err != nil {
		log.Println("删除缓存目录错误：", err)
	}
}

func CreateTmpDir() (string, error) {
	file, err := os.MkdirTemp(os.TempDir(), "too-gin")
	if err != nil {
		return file, err
	}
	TempDirPath = file
	return file, err
}

func CreateTmpFiles(name string) {
	tempDir, err := CreateTmpDir()
	if err != nil {
		return
	}
	dir, err := local.ReadDir(name)
	if err != nil {
		return
	}
	tempDir = tempDir + string(filepath.Separator)
	for _, fileInfo := range dir {
		fileName := fileInfo.Name()
		_, err := os.Stat(tempDir + fileName)
		if err == nil || os.IsExist(err) { // 如果文件存在
			_ = os.Remove(name)
		}
		file, err := local.Open(name + "/" + fileName)
		if err != nil {
			continue
		}
		bytes, err := io.ReadAll(file)
		if err != nil {
			continue
		}
		tempFile, err := os.Create(tempDir + fileName)
		if err == nil {
			tempFile.Write(bytes)
		}
		tempFile.Close()
	}
}
