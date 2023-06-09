package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

//go:embed config.schema.json
var schema string

type Files struct {
	Sources     []string `yaml:"sources" json:"sources"`
	Destination string   `yaml:"destination" json:"destination"`
	Exclude     []string `yaml:"exclude" json:"exclude"`
	Verbosity   int      `yaml:"verbosity" json:"verbosity"`
}

type DBSource struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
}

type MySQL struct {
	Sources     []DBSource `yaml:"sources" json:"sources"`
	Destination string     `yaml:"destination" json:"destination"`
}

type Telegram struct {
	Token  string `yaml:"token,omitempty" json:"token,omitempty"`
	ChatId string `yaml:"chat_id,omitempty" json:"chat_id,omitempty"`
}

type Notifications struct {
	Telegram *Telegram `yaml:"telegram,omitempty" json:"telegram,omitempty"`
}

type Config struct {
	Files         *Files        `yaml:"files,omitempty" json:"files,omitempty"`
	MySQL         *MySQL        `yaml:"mysql" json:"mysql"`
	LogDir        string        `yaml:"log_dir,omitempty" json:"log_dir,omitempty"`
	Notifications Notifications `yaml:"notifications,omitempty" json:"notifications,omitempty"`
}

func ParseFromFile(name string) (Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return Config{}, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return Config{}, err
	}

	c := Config{}
	err = yaml.Unmarshal(b, &c)

	cJSON, err := json.Marshal(c)
	if err != nil {
		return Config{}, err
	}

	res, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(schema),
		gojsonschema.NewBytesLoader(cJSON),
	)
	if err != nil {
		return Config{}, fmt.Errorf("failed to validate %s: %s", name, err)
	}
	if !res.Valid() {
		return Config{}, fmt.Errorf("failed to validate %s: %s", name, formatSchemaErrors(res))
	}

	return c, err
}

func formatSchemaErrors(res *gojsonschema.Result) string {
	return strings.Join(func() []string {
		r := make([]string, 0)
		for _, v := range res.Errors() {
			r = append(r, v.String())
		}
		return r
	}(), "; ")
}
