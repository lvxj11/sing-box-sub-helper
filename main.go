package main

import (
	"fmt"
	"os"

	//lint:ignore ST1001 我需要 dot imports 来简化调用
	. "sing-box-sub-helper/packages"
)

func main() {
	fmt.Println("Sing-Box Subscription Helper")
	fmt.Println("Version: 0.1.1")
	fmt.Println("Author: Lvxj11")
	fmt.Println("============================================================")
	fmt.Println("获取配置信息...")
	settings, err := GetSettings()
	if err != nil {
		fmt.Println("获取配置失败:", err)
		return
	}
	// 如果settings.Base64File非空且文件存在则读取
	listData := []byte{}
	if settings.Base64File != "" && FileExists(settings.Base64File) {
		fmt.Println("发现base64文件，开始读取...")
		listData, err = ReadBase64File(settings.Base64File)
		if err != nil {
			fmt.Println("读取base64文件失败:", err)
		}
	}
	if len(listData) == 0 {
		fmt.Println("从远程获取订阅数据...")
		if settings.SubscribeURL == "" || (settings.SubscribeURL[:4] != "http" && settings.SubscribeURL[:5] != "https") {
			// 如果settings.SubscribeURL为空或者不是http(s)链接则退出
			fmt.Println("订阅链接为空或不是http(s)链接，请检查配置文件！")
			return
		}
		listData, err = FetchSubscribeData(settings.SubscribeURL)
		if err != nil {
			fmt.Println("获取订阅数据失败:", err)
			return
		}
	}
	fmt.Println("写入临时节点列表文件...")
	err = os.WriteFile(settings.TempListPath, listData, 0644)
	if err != nil {
		fmt.Println("写入temp.list失败:", err)
		return
	}
	fmt.Println("开始处理订阅...")
	if err := ProcessSubscription(settings); err != nil {
		fmt.Println("处理订阅失败:", err)
	}
	// 合并到模板文件
	fmt.Println("合并模板文件...")
	if err := MergeTemplateWithSubscription(settings); err != nil {
		fmt.Println("合并模板文件失败:", err)
		return
	}
	fmt.Println("成功运行结束，输出文件：", settings.OutputPath)
}
