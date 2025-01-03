package packages

import (
	"net/url"
	"os"
)

// 验证
func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsNotExist(err) // 路径格式有效且可以访问（存在或不存在）
}
func IsValidUrl(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}

// func isFileExists(path string) bool {
// 	_, err := os.Stat(path)
// 	return err == nil
// }