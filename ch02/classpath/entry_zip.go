package classpath

import (
	"io/ioutil"
	"runtime"
	"strings"
)
import "archive/zip"
import "errors"
import "path/filepath"

const system = runtime.GOOS

type ZipEntry struct {
	absPath string
	zipRC   *zip.ReadCloser
}

func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &ZipEntry{absPath, nil}
}

func (self *ZipEntry) readClass(className string) ([]byte, Entry, error)  {
	var reader = self.zipRC
	if reader == nil {
		tempReader, err := zip.OpenReader(self.absPath)
		self.zipRC = tempReader
		reader = tempReader

		if err != nil {
			return nil, nil, err
		}
	}

	// 判断是 macos
	if system == "darwin" {
		// 统计 区分类的点 的数量, 例如 java.lang.Object.class
		// 数量为2, 转换 java.lang.Object.class --> java/lang/Object.class
		packagePoint := strings.Count(className, ".") - 1
		className = strings.Replace(className, ".", "/", packagePoint)
	}

	for _, f := range reader.File {
		if f.Name == className {
			rc, err := f.Open()
			if err != nil {
				return nil, nil, err
			}
			defer rc.Close()
			data, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, nil, err
			}
			return data, self, nil
		}
	}

	return nil, nil, errors.New("class not found: " + className)
}

func (self *ZipEntry) String() string  {
	return self.absPath
}