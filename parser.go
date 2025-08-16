package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Parser interface {
	// parse translation file
	// dir: resources directory
	// name: translation file name without extension
	Parse(dir, name string) (map[string]string, error)
}

type JsonParser struct{}

// Parse 是一个用于解析 JSON 文件的方法，它属于 JsonParser 类型。
// 该方法根据给定的目录(dir)、语言(lang)和文件名(name)来定位并解析 JSON 文件。
// 它返回一个映射，其中包含解析后的键值对，以及一个错误对象，如果解析过程中遇到任何问题，则返回相应的错误。
//
// 参数:
//   - dir: 文件所在的目录路径。
//   - name: JSON 文件的基础名，不包括扩展名。
//
// 返回值:
//   - map[string]string: 包含解析后数据的映射。
//   - error: 如果解析过程中发生错误，则返回该错误。
func (p *JsonParser) Parse(dir, name string) (map[string]string, error) {
	// 构造文件路径。
	filePath := path.Join(dir, name+".json")

	// 检查文件是否存在，如果不存在则返回空映射和nil错误。
	if !isFileExist(filePath) {
		return nil, nil
	}

	// 读取文件内容。
	byts, err := os.ReadFile(filePath)
	if err != nil {
		// 如果读取文件时发生错误，返回错误信息。
		return nil, fmt.Errorf("read file error: %v", err)
	}

	// 解析JSON数据。
	var trans map[string]string
	if err := json.Unmarshal(byts, &trans); err != nil {
		// 如果解析JSON数据时发生错误，返回错误信息。
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	// 返回解析后的映射。
	return trans, nil
}
