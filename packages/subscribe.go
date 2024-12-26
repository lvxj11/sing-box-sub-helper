package packages

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// 获取base64数据
func FetchBase64Data(settings Settings) ([]byte, error) {
	if settings.SubscribeURL == "" {
		return []byte{}, fmt.Errorf("订阅链接为空，请检查配置文件！")
	}
	resp, err := http.Get(settings.SubscribeURL)
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
func ReadBase64File(path string) ([]byte, error) {
	base64Data := []byte{}
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return []byte{}, fmt.Errorf("读取base64文件失败: %w", err)
		}
		if len(data) == 0 {
			return []byte{}, fmt.Errorf("base64文件为空")
		}
		base64Data = data
	}
	// 解码base64数据
	decodedData, err := decodeBase64(string(base64Data))
	if err != nil {
		return []byte{}, fmt.Errorf("无法解码base64数据: %w", err)
	}
	return []byte(decodedData), nil
}

func ProcessSubscription(settings Settings) error {
	input, err := os.Open(settings.TempListPath)
	if err != nil {
		return fmt.Errorf("打开输入文件失败: %w", err)
	}
	defer input.Close()

	output := []byte("[\n")
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		var jsonData []byte
		if strings.HasPrefix(line, "trojan://") {
			jsonData = processTrojan(line)
		} else if strings.HasPrefix(line, "ss://") {
			jsonData = processShadowsocks(line)
		} else {
			fmt.Println("无法识别的行格式:", line)
			continue
		}
		// 匹配tag
		if len(jsonData) != 0 && processMatchNodeJson(jsonData, settings.Filter) {
			// 将jsonData添加到output
			output = append(output, jsonData...)
			output = append(output, []byte(",\n")...)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("扫描输入文件失败: %w", err)
	}
	// 删除最后一个逗号
	output = output[:len(output)-2]
	// 添加结束符
	output = append(output, []byte("\n]")...)
	// 写入文件
	if err := os.WriteFile(settings.TempJsonPath, output, 0644); err != nil {
		return fmt.Errorf("写入输出文件失败: %w", err)
	}
	return nil
}
