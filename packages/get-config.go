package packages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetSettings 获取配置
func GetSettings() (Settings, error) {
	programDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("获取脚本目录失败:", err)
		programDir = "."
	}
	var settings Settings
	settings.ProgramDir = programDir
	settings.TemplatePath = filepath.Join(programDir, "template.json")
	settings.OutputPath = filepath.Join(programDir, "output.json")
	settings.TempListPath = filepath.Join(programDir, "temp.list")
	settings.TempJsonPath = filepath.Join(programDir, "temp.json")
	settings.SubscribeURL = ""
	settings.Base64File = "./base64.txt"
	settings.Filter = []Filter{
		{Action: "exclude", Keywords: []string{"网站|地址|剩余|过期|时间|有效|到期|官网"}},
	}
	settingsPath := filepath.Join(programDir, "settings.ini")
	if _, err := os.Stat(settingsPath); err == nil {
		fmt.Println("发现settings.ini文件，读取配置...")
		settings, err = readIniConfig(settingsPath, settings)
		if err != nil {
			fmt.Println("读取settings.ini失败:", err)
		} else {
			return settings, nil
		}
	}
	// 没有找到配置文件，返回默认配置并写入settings.ini
	err = writeIniConfig(settingsPath, settings)
	if err != nil {
		fmt.Println("写入settings.ini失败:", err)
	}
	return settings, nil
}

// 读取ini格式配置
func readIniConfig(settingsPath string, config Settings) (Settings, error) {
	settingBytes, err := os.ReadFile(settingsPath)
	if err != nil {
		fmt.Println("读取settings.ini失败:", err)
		return config, err
	}
	// 遍历每行，兼容windows和linux换行符
	lines := strings.Split(strings.ReplaceAll(string(settingBytes), "\r\n", "\n"), "\n")
	for _, line := range lines {
		// 跳过空行
		if line == "" {
			continue
		}
		// 跳过注释
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		// 解析键值对
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		switch key {
		case "subscribeURL":
			config.SubscribeURL = value
		case "templatePath":
			config.TemplatePath = value
		case "outputPath":
			config.OutputPath = value
		case "tempListPath":
			config.TempListPath = value
		case "tempJsonPath":
			config.TempJsonPath = value
		case "excludeKeywords":
			config.Filter = []Filter{
				{Action: "exclude", Keywords: []string{value}},
			}
		}
	}
	return config, nil
}

// 将默认配置写入settings.ini
func writeIniConfig(settingsPath string, config Settings) error {
	// 生成配置文件内容
	// 首行注释信息末尾加windows换行符
	settingsContent := "; 根据注释修改配置\r\n"
	settingsContent += "; subscribeURL 订阅地址\r\n"
	settingsContent += fmt.Sprintf("subscribeURL = %s\r\n", config.SubscribeURL)
	settingsContent += "; Base64File 获取到的节点列表base64加密文件\r\n"
	settingsContent += "; 如果文件路径配置并存在将跳过从远程连接获取信息\r\n"
	settingsContent += fmt.Sprintf("Base64File = %s\r\n", config.Base64File)
	settingsContent += "; templatePath 模板文件路径\r\n"
	settingsContent += fmt.Sprintf("templatePath = %s\r\n", config.TemplatePath)
	settingsContent += "; outputPath 输出文件路径\r\n"
	settingsContent += fmt.Sprintf("outputPath = %s\r\n", config.OutputPath)
	settingsContent += "; tempListPath 临时节点列表文件路径\r\n"
	settingsContent += fmt.Sprintf("tempListPath = %s\r\n", config.TempListPath)
	settingsContent += "; tempJsonPath 临时JSON文件路径\r\n"
	settingsContent += fmt.Sprintf("tempJsonPath = %s\r\n", config.TempJsonPath)
	settingsContent += "; excludeKeywords 排除关键词\r\n"
	settingsContent += fmt.Sprintf("excludeKeywords = %s\r\n", config.Filter[0].Keywords[0])
	// 写入配置文件
	if err := os.WriteFile(settingsPath, []byte(settingsContent), 0644); err != nil {
		fmt.Println("写入settings.ini失败:", err)
		return err
	}
	return nil
}
