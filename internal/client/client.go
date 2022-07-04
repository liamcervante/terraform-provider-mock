package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

type Local struct {
	ResourceDirectory string
	DataDirectory     string
}

func (local Local) ReadResource(id string, data interface{}) error {
	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return errors.New("os.ReadFile: " + err.Error())
	}

	if err := json.Unmarshal(jsonData, data); err != nil {
		return errors.New("json.Unmarshal: " + err.Error())
	}

	return nil
}

func (local Local) WriteResource(id string, data interface{}) error {
	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.New("json.MarshalIndent: " + err.Error())
	}

	if err := os.MkdirAll(local.ResourceDirectory, 0700); err != nil {
		return errors.New("os.MkdirAll: " + err.Error())
	}

	if _, err := os.Stat(jsonPath); err == nil {
		return errors.New("os.Stat: file already exists")
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return errors.New("os.WriteFile: " + err.Error())
	}

	return nil

}

func (local Local) UpdateResource(id string, data interface{}) error {
	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.New("json.MarshalIndent: " + err.Error())
	}

	if _, err := os.Stat(jsonPath); err != nil {
		if os.IsNotExist(err) {
			return errors.New("os.Stat: file does not exist")
		}
		return errors.New("os.Stat: " + err.Error())
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return errors.New("os.WriteFile: " + err.Error())
	}

	return nil

}

func (local Local) DeleteResource(id string) error {
	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))

	if err := os.Remove(jsonPath); err != nil {
		return errors.New("os.Remove: " + err.Error())
	}

	return nil
}

func (local Local) ReadDataSource(id string, data interface{}) error {
	jsonPath := path.Join(local.DataDirectory, fmt.Sprintf("%s.json", id))

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return errors.New("os.ReadFile: " + err.Error())
	}

	if err := json.Unmarshal(jsonData, data); err != nil {
		return errors.New("json.Unmarshal: " + err.Error())
	}

	return nil
}
