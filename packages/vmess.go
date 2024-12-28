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
	returnObj := Vmess{
		Tag:            jsonObj["ps"].(string),
		Type:           "vmess",
		Server:         jsonObj["add"].(string),
		ServerPort:     strToInt(jsonObj["port"].(string)),
		UUID:           jsonObj["id"].(string),
		Security:       "auto",
		AlterId:        strToInt(jsonObj["aid"].(string)),
		PacketEncoding: "xudp",
	}
	returnByte, _ := json.Marshal(returnObj)
	return returnByte
}
