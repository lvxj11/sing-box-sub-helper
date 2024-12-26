package packages

import (
	"encoding/json"
	"fmt"
	"os"
)

// 将订阅json合并到模板文件
func MergeTemplateWithSubscription(templateFile, jsonFile, outputFile string) error {
	templateBytes, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %w", err)
	}
	// 解析 JSON 数据
	var templateObj map[string]interface{}
	if err := json.Unmarshal(templateBytes, &templateObj); err != nil {
		return err
	}
	outboundsObjArr, ok := templateObj["outbounds"].([]interface{})
	if !ok {
		return fmt.Errorf("解析模板文件失败")
	}
	subJsonBytes, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("读取订阅文件失败: %w", err)
	}
	nodeTags, err := extractSubTags(subJsonBytes)
	if err != nil {
		return fmt.Errorf("读取订阅文件失败: %w", err)
	}

	outboundsMerged := []interface{}{}
	// 遍历outboundsObjArr，将每个outbound转换为map[string]interface{}
	for i, outboundObj := range outboundsObjArr {
		outbound, ok := outboundObj.(map[string]interface{})
		if !ok {
			return fmt.Errorf("解析第 %d 个 Outbound 失败", i)
		}
		// 如果outbound的type是selector或urltest则处理
		if outbound["type"] == "selector" || outbound["type"] == "urltest" {
			// 断言outbound的outbounds字段为字符串数组并遍历
			outboundTags, ok := outbound["outbounds"].([]interface{})
			if !ok {
				return fmt.Errorf("解析第 %d 个 Outbound 的 outbounds 字段失败", i)
			}

			for j, outboundTagRaw := range outboundTags {
				outboundTag, ok := outboundTagRaw.(string)
				if !ok {
					return fmt.Errorf("解析第 %d 个 Outbound 的第 %d 个 Tag 失败", i, j)
				}
				// 如果outboundTag为“{all}”则删除“{all}”并将nodeTags中的所有tag添加到outbounds
				if outboundTag == "{all}" {
					outboundTags = append(outboundTags[:j], outboundTags[j+1:]...)
					// 如果存在filter则获取filter字段
					if _, exists := outbound["filter"]; exists {
						rows := outbound["filter"].([]interface{})
						filters := extractFilters(rows)
						// 从outbound中删除filter字段
						delete(outbound, "filter")
						for _, tag := range nodeTags {
							if matchTag(tag, filters) {
								outboundTags = append(outboundTags, tag)
							}
						}
					} else {
						for _, tag := range nodeTags {
							outboundTags = append(outboundTags, tag)
						}
					}
					break
				}
			}
			outbound["outbounds"] = outboundTags
		}
		outboundsMerged = append(outboundsMerged, outbound)
	}
	// 将subJsonBytes中的节点添加到outboundsMerged
	var nodes []interface{}
	err = json.Unmarshal(subJsonBytes, &nodes)
	if err != nil {
		return fmt.Errorf("解析订阅文件失败: %w", err)
	}
	outboundsMerged = append(outboundsMerged, nodes...)
	templateObj["outbounds"] = outboundsMerged
	// 将合并后的config转换为json并保存到outputPath
	mergedBytes, err := json.MarshalIndent(templateObj, "", "  ")
	fmt.Println("合并完成，写入配置文件：", outputFile)
	if err != nil {
		return fmt.Errorf("转换合并后的配置失败: %w", err)
	}
	if err := os.WriteFile(outputFile, mergedBytes, 0644); err != nil {
		return fmt.Errorf("写入合并后的配置失败: %w", err)
	}
	return nil
}

func extractSubTags(subJsonBytes []byte) ([]string, error) {
	// 传入json字节切片，返回tag字符串数组
	var subRows []map[string]interface{}
	if err := json.Unmarshal(subJsonBytes, &subRows); err != nil {
		return nil, err
	}
	var tags []string
	for _, subJsonRaw := range subRows {
		tag, ok := subJsonRaw["tag"].(string)
		if !ok {
			continue
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func extractFilters(rows []interface{}) []Filter {
	// 传入字节切片，返回filter数组
	var filters []Filter
	for _, row := range rows {
		row, ok := row.(map[string]interface{})
		if !ok {
			continue
		}
		action, ok := row["action"].(string)
		if !ok {
			continue
		}
		keywordsRaw, ok := row["keywords"].([]interface{})
		if !ok {
			continue
		}
		var keywords []string
		for _, keywordRaw := range keywordsRaw {
			keyword, ok := keywordRaw.(string)
			if !ok {
				continue
			}
			keywords = append(keywords, keyword)
		}
		filters = append(filters, Filter{Action: action, Keywords: keywords})
	}
	return filters
}
