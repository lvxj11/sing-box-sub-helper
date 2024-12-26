package packages

import (
	"encoding/json"
	"fmt"
	"strings"
)

func processTrojan(line string) []byte {
	data := extractTrojanData(line)
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("解析trojan配置失败:", err)
		return []byte{}
	}
	return jsonData
}

func extractTrojanData(line string) TrojanConfig {
	parts := strings.Split(strings.TrimPrefix(line, "trojan://"), "@")
	password := parts[0]
	addressParts := strings.Split(parts[1], ":")
	address := addressParts[0]
	portParts := strings.Split(addressParts[1], "?")
	port := strToInt(portParts[0])
	peer := parseQueryParam(portParts[1], "peer")
	sni := parseQueryParam(portParts[1], "sni")
	tag := decodeTag(strings.Split(portParts[1], "#")[1])
	// 如果sni为空，则使用peer作为sni
	if sni == "" {
		sni = peer
	}
	tls := trojanTLS{
		Enable:     true,
		Insecure:   true,
		ServerName: sni,
	}

	return TrojanConfig{
		Tag:        tag,
		Type:       "trojan",
		Server:     address,
		ServerPort: port,
		Password:   password,
		TLS:        tls,
	}
}
