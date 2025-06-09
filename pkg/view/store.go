package view

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const ViewBaseDir = "views"

func Create(name string) (*View, error) {
	viewDir := filepath.Join(ViewBaseDir, name)
	if _, err := os.Stat(viewDir); err == nil {
		return nil, fmt.Errorf("view '%s' already exists", name)
	}
	if err := os.MkdirAll(viewDir, 0755); err != nil {
		return nil, err
	}

	v := &View{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return v.Save()
}

func Load(name string) (*View, error) {
	path := filepath.Join(ViewBaseDir, name, "view.json")
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var v View
	if err := json.Unmarshal(file, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (v *View) Save() (*View, error) {
	v.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	path := filepath.Join(ViewBaseDir, v.Name, "view.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, err
	}
	return v, nil
}

func (v *View) SetQueryFromFile(filePath string) error {
	return v.setFile("query.sql", filePath, &v.QueryFile)
}

func (v *View) SetTypeFromFile(filePath string) error {
	return v.setFile("type.graphql", filePath, &v.TypeFile)
}

func (v *View) SetTransformFromFile(filePath string) error {
	return v.setFile("transform.go", filePath, &v.TransformFile)
}

func (v *View) setFile(destName, srcPath string, ref *string) error {
	dest := filepath.Join(ViewBaseDir, v.Name, destName)
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dest, content, 0644); err != nil {
		return err
	}
	*ref = destName
	_, err = v.Save()
	return err
}
