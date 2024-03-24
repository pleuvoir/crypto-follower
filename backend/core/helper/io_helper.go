package helper

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// CloseQuietly 安静的调用Close()
func CloseQuietly(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}

// IsExists 判断文件或目录是否存在
func IsExists(p string) bool {
	_, err := os.Stat(p)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// CurrentPath 获取当前的执行文件所在的目录
func CurrentPath() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(dir), nil
}

// RootPath 获取项目根路径
func RootPath() (string, error) {
	dir, err := filepath.Abs("")
	if err != nil {
		return "", err
	}
	return filepath.Dir(dir), nil
}

// CodeRunPath 当前代码文件执行的位置
func CodeRunPath() (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename), nil
}
