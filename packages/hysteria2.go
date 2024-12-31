package packages

import (
	"encoding/json"
	"strings"
)

func processHysteria2(line string) []byte {
	parts := strings.Split(strings.TrimPrefix(line, "hysteria2://"), "@")
	password := parts[0]
	addressParts := strings.Split(parts[1], ":")
	address := addressParts[0]
	portParts := strings.Split(addressParts[1], "?")
	port := strToInt(portParts[0])
	sni := parseQueryParam(portParts[1], "sni")
	insecure := parseQueryParam(portParts[1], "insecure")
	upMbps := parseQueryParam(portParts[1], "up_mbps")
	downMbps := parseQueryParam(portParts[1], "down_mbps")
	tag := decodeTag(strings.Split(portParts[1], "#")[1])
	// 生成返回数据
	r := map[string]interface{}{}
	r["tag"] = tag
	r["type"] = "hysteria2"
	r["server"] = address
	r["server_port"] = port
	r["password"] = password
	if upMbps != "" {
		r["up_mbps"] = upMbps
	}
	if downMbps != "" {
		r["down_mbps"] = downMbps
	}
	// 生成返回数据tls字段
	r_tls := map[string]interface{}{}
	r_tls["enabled"] = true
	r_tls["insecure"] = insecure == "1"
	r_tls["alpn"] = []string{"h3"}
	if sni != "" {
		r_tls["server_name"] = sni
	}
	r["tls"] = r_tls
	// 返回数据编码为json
	r_json, err := json.Marshal(r)
	if err != nil {
		PrintRed("hysteria2配置json编码失败:" + err.Error())
		return []byte{}
	}

	return r_json
}
