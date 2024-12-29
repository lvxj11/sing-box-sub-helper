package packages

import (
	"encoding/json"
	"fmt"
	"strings"
)

func processTrojan(line string) []byte {
	parts := strings.Split(strings.TrimPrefix(line, "trojan://"), "@")
	password := parts[0]
	addressParts := strings.Split(parts[1], ":")
	address := addressParts[0]
	portParts := strings.Split(addressParts[1], "?")
	port := strToInt(portParts[0])
	insecure := parseQueryParam(portParts[1], "allowInsecure")
	sni := parseQueryParam(portParts[1], "sni")
	peer := parseQueryParam(portParts[1], "peer")
	tag := decodeTag(strings.Split(portParts[1], "#")[1])
	// 如果sni为空，则使用peer作为sni
	if sni == "" {
		sni = peer
	}
	// 生成返回数据
	r := map[string]interface{}{}
	r["tag"] = tag
	r["type"] = "trojan"
	r["server"] = address
	r["server_port"] = port
	r["password"] = password
	// 生成返回数据tls字段
	r_tls := map[string]interface{}{}
	r_tls["enabled"] = true
	r_tls["insecure"] = insecure == "1"
	if sni != "" {
		r_tls["server_name"] = sni
	}
	r["tls"] = r_tls
	// 返回数据编码为json
	r_json, err := json.Marshal(r)
	if err != nil {
		fmt.Println("trojan配置json编码失败:", err)
		return []byte{}
	}

	return r_json
}
