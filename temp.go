package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var TempDirPath string

// destroyTemp 用于删除指定路径下的所有文件和文件夹。
// 该函数设计为仅在路径包含 "temp" 字符串时执行操作，以防止误删其他重要数据。
// 参数：
//
//	path - 要删除的文件或目录的路径。
//
// 返回值：
//
//	error - 如果删除过程中发生错误，返回相应的错误；否则返回 nil。
func destroyTemp(path string) error {
	// 只有当路径中包含 "temp" 字符串时才执行删除操作（目前此判断未实现具体逻辑）。
	if strings.Contains(path, "temp") {

	}

	// 使用 filepath.Walk 遍历指定路径下的所有文件和子目录。
	// 注意：也可以使用更高效的 filepath.WalkDir 函数，但此处选择使用 filepath.Walk。
	err := filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		// 如果文件信息为 nil，说明无法获取该文件的信息，直接返回错误。
		if nil == fi {
			return err
		}

		// 如果当前项不是目录，则将其作为普通文件删除。
		if !fi.IsDir() {
			err := os.Remove(path)
			if err != nil {
				return err // 删除文件失败时返回错误。
			}
			return nil
		}

		// 如果当前项是目录，则递归删除整个目录及其内容。
		err = os.RemoveAll(path)
		if err != nil {
			return err // 删除目录失败时返回错误。
		}
		return nil
	})

	// 返回遍历和删除过程中的最终结果（nil 表示成功，非 nil 表示出错）。
	return err
}

// DestroyTempDir 删除临时目录及其所有内容。
// 该函数旨在清理不再需要的临时文件，以释放系统资源。
// 没有输入参数，也不返回任何值。
// 当无法删除目录时，会记录一个错误日志。
func DestroyTempDir() {
	err := os.RemoveAll(TempDirPath)
	if err != nil {
		log.Println("删除缓存目录错误：", err)
	}
}

// CreateTmpDir 创建一个临时目录用于too-gin框架。
// 该函数不接受任何参数。
// 返回值是一个字符串，表示新创建的临时目录的路径，和一个错误值，如果创建过程中发生错误。
func CreateTmpDir() (string, error) {
	// 使用os.MkdirTemp在系统临时目录下创建一个以"too-gin"为前缀的临时目录。
	file, err := os.MkdirTemp(os.TempDir(), "too-gin")
	if err != nil {
		// 如果创建临时目录时发生错误，返回空字符串和错误详情。
		return file, err
	}
	// 将新创建的临时目录路径赋值给全局变量TempDirPath，以便后续使用。
	TempDirPath = file
	// 返回新创建的临时目录路径和nil错误，表示操作成功。
	return file, err
}

// CreateTmpFiles 创建临时文件。
// 该函数根据给定的名称从一个目录中读取所有文件，并将它们复制到一个临时目录中。
// 参数:
//
//	name - 源目录的名称，从该目录中读取文件。
func CreateTmpFiles(name string) {
	// 创建一个临时目录。
	tempDir, err := CreateTmpDir()
	if err != nil {
		return
	}

	// 读取源目录中的文件信息。
	dir, err := local.ReadDir(name)
	if err != nil {
		return
	}

	// 确保临时目录路径以路径分隔符结尾。
	tempDir = tempDir + string(filepath.Separator)

	// 遍历源目录中的所有文件信息。
	for _, fileInfo := range dir {
		fileName := fileInfo.Name()

		// 检查临时目录中是否已存在同名文件。
		_, err := os.Stat(tempDir + fileName)
		if err == nil || os.IsExist(err) { // 如果文件存在
			// 尝试删除源文件，如果失败则忽略错误。
			_ = os.Remove(name)
		}

		// 打开源文件。
		file, err := local.Open(name + "/" + fileName)
		if err != nil {
			continue
		}

		// 读取源文件的全部内容。
		bytes, err := io.ReadAll(file)
		if err != nil {
			continue
		}

		// 在临时目录中创建新文件。
		tempFile, err := os.Create(tempDir + fileName)
		if err == nil {
			// 将源文件内容写入新文件。
			_, err := tempFile.Write(bytes)
			if err != nil {
				return
			}
		}

		// 关闭新创建的文件。
		err = tempFile.Close()
		if err != nil {
			return
		}
	}
}
