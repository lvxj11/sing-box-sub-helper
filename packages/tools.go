package packages

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parsePort(portStr string) int {
	port, _ := strconv.Atoi(portStr)
	return port
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func parseQueryParam(query, key string) string {
	params := strings.Split(query, "#")[0]
	for _, param := range strings.Split(params, "&") {
		kv := strings.Split(param, "=")
		if kv[0] == key {
			return kv[1]
		}
	}
	return ""
}

func decodeBase64(encoded string) (string, error) {
	// str如果长度不是 4 的倍数，补充填充字符 '='
	padding := len(encoded) % 4
	if padding != 0 {
		encoded += strings.Repeat("=", 4-padding)
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	return string(decoded), err
}

func decodeTag(encodedTag string) string {
	tag, _ := url.QueryUnescape(encodedTag)
	return tag
}

// 传入字节切片和filters，将切片中的tag节点解码并匹配
func processMatchNodeJson(jsonData []byte, filters []Filter) bool {
	// 解析tag节点
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return false
	}
	tag, ok := data["tag"].(string)
	if !ok {
		return false
	}
	// 匹配tag
	return matchTag(tag, filters)
}

// 传入字符串和Filters，返回是否匹配
func matchTag(tag string, filters []Filter) bool {
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

// 传入subJsonRaw，返回tag字符串数组
func extractSubTags(subJsonBytes []byte) ([]string, error) {
	var subRows []map[string]interface{}
	if err := json.Unmarshal(subJsonBytes, &subRows); err != nil {
		return nil, err
	}
	var tags []string
	for _, subJsonRaw := range subRows {
		tag, ok := subJsonRaw["tag"].(string)
		if !ok {
			continue
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// 传入字节切片，返回filter数组
func extractFilters(rows []interface{}) []Filter {
	var filters []Filter
	for _, row := range rows {
		row, ok := row.(map[string]interface{})
		if !ok {
			continue
		}
		action, ok := row["action"].(string)
		if !ok {
			continue
		}
		keywordsRaw, ok := row["keywords"].([]interface{})
		if !ok {
			continue
		}
		var keywords []string
		for _, keywordRaw := range keywordsRaw {
			keyword, ok := keywordRaw.(string)
			if !ok {
				continue
			}
			keywords = append(keywords, keyword)
		}
		filters = append(filters, Filter{Action: action, Keywords: keywords})
	}
	return filters
}
