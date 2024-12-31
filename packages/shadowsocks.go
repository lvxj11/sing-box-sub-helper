package packages

import (
	"encoding/json"
	"strings"
)

func processShadowsocks(line string) []byte {
	encoded := strings.Split(strings.TrimPrefix(line, "ss://"), "@")[0]
	decodedConfig, err := decodeBase64([]byte(encoded))
	if err != nil {
		PrintRed("shadowsocks配置base64解密失败:" + err.Error())
		return []byte{}
	}
	parts := strings.Split(string(decodedConfig), ":")
	method := parts[0]
	password := parts[1]
	addressParts := strings.Split(strings.Split(strings.TrimPrefix(line, "ss://"), "@")[1], ":")
	address := addressParts[0]
	portParts := strings.Split(addressParts[1], "#")
	port := strToInt(portParts[0])
	tag := decodeTag(portParts[1])

	// 生成返回数据
	r := map[string]interface{}{}
	r["tag"] = tag
	r["type"] = "shadowsocks"
	r["server"] = address
	r["server_port"] = port
	r["method"] = method
	r["password"] = password

	// 返回数据编码为json
	r_json, err := json.Marshal(r)
	if err != nil {
		PrintRed("shadowsocks配置json编码失败:" + err.Error())
		return []byte{}
	}

	return r_json
}
