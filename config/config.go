package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"spider/cookie"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Database []MySQLConfig `yaml:"database"`
	Tasks    []TaskConfig  `yaml:"tasks"`
	Includes []string      `yaml:"includes"`
}

type MySQLConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// TaskConfig 是任务的结构体定义
type TaskConfig struct {
	Name        string        `yaml:"name"`
	Schedule    string        `yaml:"schedule"`
	HandlerName string        `yaml:"handler_name"`
	Cookie      CookieConfig  `yaml:"cookie"`
	Request     RequestConfig `yaml:"request"`
	Storage     StorageConfig `yaml:"storage"`
}

// CookieConfig 是cookie配置的结构定义
type CookieConfig struct {
	FetcherName string          `yaml:"fetcher_name"`
	LoginURL    string          `yaml:"login_url"`
	Actions     []cookie.Action `yaml:"actions"`
	Schedule    string          `yaml:"schedule"`
}

type RequestConfig struct {
	URL        string            `yaml:"url"`
	Method     string            `yaml:"method"`
	Headers    map[string]string `yaml:"headers"`
	Body       string            `yaml:"body"`
	BodyParams map[string]string `yaml:"body_params"`
	Params     map[string]string `yaml:"params"`
}

func (r *RequestConfig) BuildURL() string {
	u, err := url.Parse(r.URL)
	if err != nil {
		return r.URL // 错误时返回原始 URL
	}

	q := u.Query()
	for key, value := range r.Params {
		q.Add(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (r *RequestConfig) BuildBody() string {
	body := r.Body
	for key, value := range r.BodyParams {
		body = strings.Replace(body, fmt.Sprintf("{%s}", key), value, -1)
	}
	return body
}

type StorageConfig struct {
	StorageName string       `yaml:"storage_name"`
	MySQLConfig *MySQLConfig `yaml:"mysql_config"`
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

	config.AssignDatabaseConfigs()

	// 加载所有包含的子配置文件
	for _, include := range config.Includes {
		subConfig, err := loadSubConfig(include)
		if err != nil {
			return nil, err
		}
		config.MergeTasks(subConfig.Tasks)
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

// AssignDatabaseConfigs 是根据配置的存储器名称加载存储配置
func (c *Config) AssignDatabaseConfigs() {
	for i := range c.Tasks {
		task := &c.Tasks[i]
		if task.Storage.MySQLConfig == nil {
			for _, db := range c.Database {
				if db.Name == task.Storage.StorageName {
					task.Storage.MySQLConfig = &db
				}
			}
		}
	}
}

// MergeTasks 是将includes文件里的任务合并
func (c *Config) MergeTasks(subTasks []TaskConfig) {
	taskMap := make(map[string]bool)
	for _, task := range c.Tasks {
		taskMap[task.Name] = true
	}

	for _, subTask := range subTasks {
		if !taskMap[subTask.Name] {
			c.Tasks = append(c.Tasks, subTask)
		}
	}
}
