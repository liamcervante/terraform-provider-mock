package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/liamcervante/terraform-provider-mock/internal/values"
)

type Local struct {
	ResourceDirectory string
	DataDirectory     string
}

func (local Local) ReadResource(ctx context.Context, id string) (*values.Resource, error) {
	tflog.Trace(ctx, "Local.ReadResource")

	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.New("os.ReadFile: " + err.Error())
	}

	var value values.Resource
	if err := json.Unmarshal(jsonData, &value); err != nil {
		return nil, errors.New("json.Unmarshal: " + err.Error())
	}

	return &value, nil
}

func (local Local) WriteResource(ctx context.Context, id string, value *values.Resource) error {
	tflog.Trace(ctx, "Local.WriteResource")

	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))

	jsonData, err := json.MarshalIndent(value, "", "  ")
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

func (local Local) UpdateResource(ctx context.Context, id string, value *values.Resource) error {
	tflog.Trace(ctx, "Local.UpdateResource")

	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))
	jsonData, err := json.MarshalIndent(value, "", "  ")
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

func (local Local) DeleteResource(ctx context.Context, id string) error {
	tflog.Trace(ctx, "Local.DeleteResource")

	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))

	if err := os.Remove(jsonPath); err != nil {
		return errors.New("os.Remove: " + err.Error())
	}

	return nil
}

func (local Local) ReadDataSource(ctx context.Context, id string) (*values.Resource, error) {
	tflog.Trace(ctx, "Local.ReadDataSource")

	jsonPath := path.Join(local.DataDirectory, fmt.Sprintf("%s.json", id))

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, errors.New("os.ReadFile: " + err.Error())
	}

	var value values.Resource
	if err := json.Unmarshal(jsonData, &value); err != nil {
		return nil, errors.New("json.Unmarshal: " + err.Error())
	}

	return &value, nil
}
