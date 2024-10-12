package config

import (
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Tasks    []TaskConfig `yaml:"tasks"`
	Includes []string     `yaml:"includes"`
}

// TaskConfig 是任务的结构体定义
type TaskConfig struct {
	Name     string        `yaml:"name"`
	Schedule string        `yaml:"schedule"`
	Cookie   CookieConfig  `yaml:"cookie"`
	Request  RequestConfig `yaml:"request"`
	Storage  StorageConfig `yaml:"storage"`
}

type CookieConfig struct {
	Method         string `yaml:"method"`
	UpdateInterval string `yaml:"update_interval"`
}

type RequestConfig struct {
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

type StorageConfig struct {
	Type       string `yaml:"type"`
	Connection string `yaml:"connection"`
	Path       string `yaml:"path"`
	ReplaceOld bool   `yaml:"replace_old"`
}

// LoadConfig 加载主配置文件，并读取包含的所有子任务配置文件
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// 加载所有包含的子配置文件
	for _, include := range config.Includes {
		subConfig, err := loadSubConfig(include)
		if err != nil {
			return nil, err
		}
		config.Tasks = append(config.Tasks, subConfig.Tasks...)
	}

	return &config, nil
}

// loadSubConfig 加载子任务配置文件
func loadSubConfig(filePath string) (*Config, error) {
	fullPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	var subConfig Config
	if err := yaml.Unmarshal(data, &subConfig); err != nil {
		return nil, err
	}

	return &subConfig, nil
}
