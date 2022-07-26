package dynamic

import (
	"encoding/json"
	"os"
)

type Reader interface {
	Read() (map[string]Attribute, error)
}

type FileReader struct {
	File string
}

type StringReader struct {
	Data string
}

func (r FileReader) Read() (map[string]Attribute, error) {
	data, err := os.ReadFile("dynamic_resources.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var dynamicResources map[string]Attribute
	if len(data) > 0 {
		if err := json.Unmarshal(data, &dynamicResources); err != nil {
			return nil, err
		}
	}

	return dynamicResources, nil
}

func (r StringReader) Read() (map[string]Attribute, error) {
	var dynamicResources map[string]Attribute
	if len(r.Data) > 0 {
		if err := json.Unmarshal([]byte(r.Data), &dynamicResources); err != nil {
			return nil, err
		}
	}
	return dynamicResources, nil
}
