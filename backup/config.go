package backup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var Conf *Config

type Config struct {
	AccessKey string   `json:"access_key"`
	Secretkey string   `json:"secret_key"`
	Domain    string   `json:"domain"`
	Bucket    string   `json:"bucket"`
	Path      string   `json:"path"`
	Ignories  []string `json:"ignories"`
}

func (c *Config) Validate() bool {
	if c.AccessKey == "" {
		return false
	}

	if c.Secretkey == "" {
		return false
	}

	if c.Bucket == "" {
		return false
	}

	if c.Domain == "" {
		return false
	}

	return true
}

func (c *Config) Restore(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0700)
}

func NewConfig() *Config {
	return &Config{
		Path: fmt.Sprintf("%s/QBackup", os.Getenv("HOME")),
		Ignories: []string{
			".DS_Store",
		},
	}
}
