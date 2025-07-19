package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Parser interface {
	// parse translation file
	// dir: translation file directory
	// name: translation file name without extension
	Parse(dir string, name string) (map[string]string, error)
}

type JsonParser struct{}

func (p *JsonParser) Parse(dir string, name string) (map[string]string, error) {
	filePath := path.Join(dir, name+".json")
	if !isFileExist(filePath) {
		return nil, nil
	}

	byts, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file error: %v", err)
	}
	var trans map[string]string
	if err := json.Unmarshal(byts, &trans); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return trans, nil
}
