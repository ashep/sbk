package config

import (
	_ "embed"
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
	Token  string `yaml:"token" json:"token"`
	ChatId string `yaml:"chat_id" json:"chat_id"`
}

type Notifications struct {
	Telegram Telegram `yaml:"telegram" json:"telegram"`
}

type Config struct {
	Files         Files         `yaml:"files" json:"files"`
	MySQL         MySQL         `yaml:"mysql" json:"mysql"`
	Notifications Notifications `yaml:"notifications" json:"notifications"`
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

	c := Config{Files: Files{
		Verbosity: 5,
	}}
	err = yaml.Unmarshal(b, &c)

	res, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(schema),
		gojsonschema.NewGoLoader(c),
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
