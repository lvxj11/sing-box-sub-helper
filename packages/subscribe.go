package packages

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// 获取base64数据
func FetchBase64Data(u string) ([]byte, error) {
	if u == "" {
		return []byte{}, fmt.Errorf("订阅链接为空，请检查配置文件！")
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		PrintRed("建立请求出错...")
		return []byte{}, err
	}
	req.Header.Set("User-Agent", "curl")
	client := &http.Client{}
	resp, err := client.Do(req)
	// resp, err := http.Get(u)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("意外状态代码: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// 读取base64文件并解码
func ReadBase64FileDecode(path string) ([]byte, error) {
	base64Data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("读取base64文件失败: %w", err)
	}
	if len(base64Data) == 0 {
		return []byte{}, fmt.Errorf("base64文件为空")
	}
	// 解码base64数据
	decodedData, err := decodeBase64(base64Data)
	if err != nil {
		return []byte{}, fmt.Errorf("无法解码base64数据: %w", err)
	}
	return decodedData, nil
}
func ConvertSubscriptionToJson(path string, filter []Filter) ([]byte, error) {
	// 订阅连接列表文件转换为json格式
	subList, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("读取订阅列表文件失败: %w", err)
	}
	if len(subList) == 0 {
		return []byte{}, fmt.Errorf("订阅列表为空")
	}
	// 将可能的windows或linux换行符转换为unix换行符
	subListStr := strings.ReplaceAll(string(subList), "\r\n", "\n")
	// 按行遍历订阅连接列表
	output := []byte("[\n")
	for _, line := range strings.Split(subListStr, "\n") {
		// 跳过空行和注释行
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var jsonData []byte
		// 获取前缀
		prefix := strings.Split(line, "://")[0]
		switch prefix {
		case "trojan":
			jsonData = processTrojan(line)
		case "ss":
			jsonData = processShadowsocks(line)
		case "vmess":
			jsonData = processVmess(line)
		case "hysteria2":
			jsonData = processHysteria2(line)
		case "ssr":
			fmt.Println("不支持的节点协议：ssr")
		default:
			fmt.Println("无法识别的行格式:", line)
			continue
		}
		// 匹配tag
		if len(jsonData) != 0 && processMatchNodeJson(jsonData, filter) {
			// 将jsonData添加到output
			output = append(output, jsonData...)
			output = append(output, []byte(",\n")...)
		}
	}
	// 删除最后一个逗号
	output = output[:len(output)-2]
	// 添加结束符
	output = append(output, []byte("\n]")...)
	return output, nil
}

func decodeTag(encodedTag string) string {
	// 解码节点标签，url请求字符串（百分号字符）解码
	tag, _ := url.QueryUnescape(encodedTag)
	return tag
}

func parseQueryParam(query, key string) string {
	// 解析query参数
	params := strings.Split(query, "#")[0]
	for _, param := range strings.Split(params, "&") {
		kv := strings.Split(param, "=")
		if kv[0] == key {
			return kv[1]
		}
	}
	return ""
}

func processMatchNodeJson(jsonData []byte, filters []Filter) bool {
	// 传入字节切片和filters，将切片中的tag节点解码并匹配
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
