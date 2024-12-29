package packages

import (
	"encoding/json"
	"fmt"
	"strings"
)

func processVmess(line string) []byte {
	encodedData := strings.TrimPrefix(line, "vmess://")
	decodedData, err := decodeBase64([]byte(encodedData))
	if err != nil {
		fmt.Println("vmess配置base64解密失败:", err)
		return []byte{}
	}
	jsonObj := make(map[string]interface{})
	err = json.Unmarshal(decodedData, &jsonObj)
	if err != nil {
		fmt.Println("vmess配置json解码失败:", err)
		return []byte{}
	}

	// 生成返回数据
	r := map[string]interface{}{}
	r["tag"] = jsonObj["ps"]
	r["type"] = "vmess"
	r["server"] = jsonObj["add"]
	r["server_port"] = strToInt(jsonObj["port"].(string))
	r["uuid"] = jsonObj["id"]
	r["security"] = "auto"
	r["alter_id"] = strToInt(jsonObj["aid"].(string))
	r["packet_encoding"] = "xudp"

	// 返回数据编码为json
	r_json, err := json.Marshal(r)
	if err != nil {
		fmt.Println("vmess配置json编码失败:", err)
		return []byte{}
	}

	return r_json
}
