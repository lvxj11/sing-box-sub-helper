package packages

import (
	"encoding/json"
	"fmt"
	"strings"
)

func processShadowsocks(line string) []byte {
	data := extractShadowsocksData(line)
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("解析shadowsocks配置失败:", err)
		return []byte{}
	}
	return jsonData
}

func extractShadowsocksData(line string) ShadowsocksConfig {
	encoded := strings.Split(strings.TrimPrefix(line, "ss://"), "@")[0]
	decodedConfig, err := decodeBase64(encoded)
	if err != nil {
		fmt.Println("解码shadowsocks配置失败:", err)
		return ShadowsocksConfig{}
	}
	parts := strings.Split(string(decodedConfig), ":")
	method := parts[0]
	password := parts[1]
	addressParts := strings.Split(strings.Split(strings.TrimPrefix(line, "ss://"), "@")[1], ":")
	address := addressParts[0]
	portParts := strings.Split(addressParts[1], "#")
	port := parsePort(portParts[0])
	tag := decodeTag(portParts[1])

	return ShadowsocksConfig{
		Tag:        tag,
		Type:       "shadowsocks",
		Server:     address,
		ServerPort: port,
		Method:     method,
		Password:   password,
	}
}
