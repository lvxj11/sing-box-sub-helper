package packages

import (
	"encoding/base64"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func SeveFile(path string, content []byte) error {
	// 保存文件
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func strToInt(str string) int {
	// 将字符串转换为整数，转换失败返回0
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

func decodeBase64(encoded []byte) ([]byte, error) {
	// 解码base64数据
	// str如果长度不是 4 的倍数，补充填充字符 '='
	padding := len(encoded) % 4
	if padding != 0 {
		encoded = append(encoded, []byte(strings.Repeat("=", 4-padding))...)
	}
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
	n, err := base64.StdEncoding.Decode(decoded, encoded)
	return decoded[:n], err
}

func matchTag(tag string, filters []Filter) bool {
	// 匹配节点标签
	// 定义默认返回值, 默认为true
	// 同filter中keywerds为或关系
	// 不同filter为与关系
	defRet := true
	for _, filter := range filters {
		for _, keyword := range filter.Keywords {
			re, err := regexp.Compile(keyword)
			if err != nil {
				continue
			}
			if filter.Action == "include" {
				if re.MatchString(tag) {
					defRet = true
					break
				} else {
					defRet = false
				}
			}
			if filter.Action == "exclude" {
				if re.MatchString(tag) {
					defRet = false
					break
				} else {
					defRet = true
				}
			}
		}
		if !defRet {
			return defRet
		}
	}
	return defRet
}
