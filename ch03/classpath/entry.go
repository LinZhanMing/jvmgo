package classpath

import "os"
import "strings"

const pathListSeparator = string(os.PathListSeparator) // 路径分隔符
type Entry interface { // 定义 Entry 接口
	readClass(className string) ([]byte, Entry, error) // 负责寻找和加载class文件
	String() string // 相当于 toString
}

// 创建 Entry 实例
func newEntry(path string) Entry {
	// /aaa/pathA:/bbb/pathB:/ccc/pathC
	if strings.Contains(path, pathListSeparator) {
		return newCompositeEntry(path)
	}
	if strings.Contains(path, "*") {
		return newWildcardEntry(path)
	}
	if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".ZIP") {
		return newZipEntry(path)
	}
	return newDirEntry(path)
}
