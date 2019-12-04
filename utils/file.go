package utils

import (
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 目录下所有的文件夹
func GetDirList(dirPath string) ([]string, error) {
	var dirList []string
	err := filepath.Walk(dirPath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				dirList = append(dirList, path)
				return nil
			}

			return nil
		})
	return dirList, err
}

// 获取一个目录下所有文件信息，包含子目录
func GetDirListAll(files []os.FileInfo, dirPath string) ([]os.FileInfo, error) {
	err := filepath.Walk(dirPath, func(dPath string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			files = append(files, f)
		} else {
			GetDirListAll(files, strings.ReplaceAll(filepath.Join(dPath, f.Name()), "\\", "/"))
		}
		return nil
	})
	return files, err
}

// 获取当前路径下所有文件
// ioutil中提供了一个非常方便的函数函数ReadDir，
// 他读取目录并返回排好序的文件以及子目录名([]os.FileInfo)
func GetFileList(path string) ([]os.FileInfo, error) {
	readerInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if readerInfos == nil {
		return nil, nil
	}
	return readerInfos, nil
}

// 判断路径是否为目录
func IsExistDir(dirPath string) bool {
	if IsStringEmpty(dirPath) {
		return false
	}
	info, err := os.Stat(dirPath)
	if err != nil || !os.IsExist(err) || !info.IsDir() {
		return false
	}
	return true
}

// 判断文件是否存在：存在，返回true，否则返回false
func IsFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil || os.IsNotExist(err) || info.IsDir() {
		return false
	}
	return true
}

// 判断所给路径文件/文件夹是否存在
func IsExists(path string) bool {
	if IsStringEmpty(path) {
		return false
	}
	// os.Stat获取文件信息
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径文件/文件夹是否不存在
func IsNotExists(path string) bool {
	return !IsExists(path)
}

// 获取当前程序运行所在路径
func OsPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// 获取路径中的文件的后缀
func GetSuffix(filePath string) string {
	ext := filepath.Ext(filePath)
	return ext
}

// 获取路径中的目录及文件名
func GetDirFile(filePath string) (dir, file string) {
	paths, fileName := filepath.Split(filePath)
	return paths, fileName
}

// 获取父级目录
func ParentDirectory(dir string) string {
	return filepath.Join(dir, "..")
}

// 目录分隔符转换
func PathSeparatorSlash(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// 获取上下文路径，传入指定目录截取前一部分
func ContextPath(root string) (path string, err error) {
	// 获取当前绝对路径
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	index := strings.LastIndex(dir, root)
	if len(dir) < len(root) || index <= 0 {
		return dir, errors.New("错误：路径不正确")
	}
	return dir[0 : index+len(root)], nil
}

// 创建所有不存在的层级目录
func Mkdir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0711)
		return err
	}
	return nil
}

// 创建文件
func CreateFile(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		_, err = os.Create(filePath)
		return err
	}
	return nil
}

// 获取文件MIME类型
// 见函数http.ServeContent
func GetContentType(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	ctype := mime.TypeByExtension(filepath.Ext(fi.Name()))
	if ctype == "" {
		// read a chunk to decide between utf-8 text and binary
		var buf [512]byte
		n, _ := io.ReadFull(f, buf[:])
		// 根据前512个字节的数据判断MIME类型
		ctype = http.DetectContentType(buf[:n])
		_, err := f.Seek(0, io.SeekStart) // rewind to output whole file
		if err != nil {
			return "", err
		}
	}
	return ctype, nil
}
