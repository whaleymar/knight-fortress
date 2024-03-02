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

func StructToYaml(data interface{}) (string, error) {
	b, err := yaml.Marshal(&data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func LoadStruct(path string, dst interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, dst)
}

func YamlToStruct(ymlString string, dst interface{}) error {
	b := []byte(ymlString)
	return yaml.Unmarshal(b, dst)
}
