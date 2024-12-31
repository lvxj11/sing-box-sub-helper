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
		PrintRed("获取脚本目录失败:" + err.Error())
		programDir = "."
	}
	var settings Settings
	settings.ProgramDir = programDir
	settings.SubscribeURL = ""
	settings.Base64File = "./base64.txt"
	settings.TempListPath = filepath.Join(programDir, "temp.list")
	settings.TempJsonPath = filepath.Join(programDir, "temp.json")
	settings.TemplatePath = filepath.Join(programDir, "template.json")
	settings.OutputPath = filepath.Join(programDir, "output.json")
	settings.StartStep = 0

	settings.Filter = []Filter{
		{Action: "exclude", Keywords: []string{"网站|地址|剩余|过期|时间|有效|到期|官网"}},
	}
	settingsPath := filepath.Join(programDir, "settings.ini")
	if _, err := os.Stat(settingsPath); err == nil {
		fmt.Println("发现settings.ini文件，读取配置...")
		settings, err = readIniConfig(settingsPath, settings)
		if err != nil {
			PrintRed("读取settings.ini失败:" + err.Error())
		} else {
			return settings, nil
		}
	}
	// 没有找到配置文件，返回默认配置并写入settings.ini
	err = writeIniConfig(settingsPath, settings)
	if err != nil {
		PrintRed("写入settings.ini失败:" + err.Error())
	}
	return settings, nil
}

// 读取ini格式配置
func readIniConfig(settingsPath string, settings Settings) (Settings, error) {
	settingBytes, err := os.ReadFile(settingsPath)
	if err != nil {
		PrintRed("读取settings.ini失败:" + err.Error())
		return settings, err
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
		keyArr := strings.SplitN(key, "-", 2)
		switch keyArr[0] {
		case "subscribeURL":
			if isValidUrl(value) {
				settings.SubscribeURL = value
				settings.SubURLs = append(settings.SubURLs, SubURL{Tag: keyArr[1], URL: value})
			} else {
				return settings, fmt.Errorf("订阅地址验证错误：%s", value)
			}
		case "base64File":
			if isValidPath(value) {
				settings.Base64File = value
			}
		case "tempListPath":
			if isValidPath(value) {
				settings.TempListPath = value
			}
		case "tempJsonPath":
			if isValidPath(value) {
				settings.TempJsonPath = value
			}
		case "templatePath":
			if isValidPath(value) {
				settings.TemplatePath = value
			} else {
				return settings, fmt.Errorf("模板文件路径验证错误：%s", value)
			}
		case "excludeKeywords":
			settings.Filter = []Filter{{Action: "exclude", Keywords: []string{value}}}
		case "outputPath":
			if isValidPath(value) {
				settings.OutputPath = value
			}
		case "startStep":
			settings.StartStep = strToInt(value)
			if settings.StartStep != 0 {
				// 使用 ANSI 转义码显示红色并加粗
				fmt.Println("\033[1;31m检测到修改了默认步骤，请确保相应步骤所需数据存在并正确配置路径。\033[0m")
				if settings.StartStep > 4 {
					fmt.Println("检测到开始步骤大于4，将自动重置为4。")
					settings.StartStep = 4
				}
				fmt.Println("开始步骤为：", settings.StartStep)
			}
		}
	}
	return settings, nil
}

// 将默认配置写入settings.ini
func writeIniConfig(settingsPath string, config Settings) error {
	// 生成配置文件内容
	// 首行注释信息末尾加windows换行符
	settingsContent := "; 根据注释修改配置，初始为默认设置，如设置值验证失败将采用默认值。\r\n"
	settingsContent += "; 必填项（除非修改了开始步骤）：subscribeURL 订阅地址\r\n"
	settingsContent += fmt.Sprintf("subscribeURL = %s\r\n", config.SubscribeURL)
	settingsContent += ";必填项：templatePath 模板文件路径\r\n"
	settingsContent += fmt.Sprintf("templatePath = %s\r\n", config.TemplatePath)
	settingsContent += ";outputPath 输出文件路径，默认值为程序目录下“output.json”\r\n"
	settingsContent += fmt.Sprintf("outputPath = %s\r\n", config.OutputPath)
	settingsContent += "; base64File 获取到的节点列表base64加密文件\r\n"
	settingsContent += "; 如果文件路径配置并存在将跳过从远程连接获取信息\r\n"
	settingsContent += fmt.Sprintf("Base64File = %s\r\n", config.Base64File)
	settingsContent += "; tempListPath 临时节点列表文件路径\r\n"
	settingsContent += fmt.Sprintf("tempListPath = %s\r\n", config.TempListPath)
	settingsContent += "; tempJsonPath 临时JSON文件路径\r\n"
	settingsContent += fmt.Sprintf("tempJsonPath = %s\r\n", config.TempJsonPath)
	settingsContent += "; excludeKeywords 排除关键词\r\n"
	settingsContent += fmt.Sprintf("excludeKeywords = %s\r\n", config.Filter[0].Keywords[0])
	settingsContent += "; startStep 开始步骤,1.远程订阅地址，2.读取base64File，3.读取临时列表，\r\n"
	settingsContent += "; 4.读取临时json，后续步骤由于会缺少数据所以大于等于4的步骤都会从第4步开始。\r\n"
	settingsContent += "; 修改步骤后请确保相应数据文件路径配置正确，否则会因为无数据出错。\r\n"
	settingsContent += fmt.Sprintf("startStep = %d\r\n", config.StartStep)
	// 写入配置文件
	if err := os.WriteFile(settingsPath, []byte(settingsContent), 0644); err != nil {
		PrintRed("写入settings.ini失败:" + err.Error())
		return err
	}
	return nil
}
