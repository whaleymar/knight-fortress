package sys

import (
	"os"

	"gopkg.in/yaml.v3"
)

func SaveStruct(path string, data interface{}) error {
	b, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func LoadStruct(path string, dst interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, dst)
}
