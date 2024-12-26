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
	// 如果settings.StartStep小于等于1，从远程获取开始
	if settings.StartStep <= 1 {
		fmt.Println("从远程获取订阅数据...")
		base64Data, err := FetchBase64Data(settings)
		if err != nil {
			fmt.Println("获取订阅数据失败:", err)
			return
		}
		err = SeveFile(settings.Base64File, base64Data)
		if err != nil {
			fmt.Println("保存base64文件失败:", err)
			return
		}
	}
	// 如果settings.StartStep小于等于2
	if settings.StartStep <= 2 {
		fmt.Println("发现base64文件，开始读取...")
		listData, err := ReadBase64File(settings.Base64File)
		if err != nil {
			fmt.Println("读取base64文件失败:", err)
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
