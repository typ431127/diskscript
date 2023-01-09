package diskscript

import "net/url"

var (
	AlertmanagerURL *url.URL
	Hostname        string
	Conf            *YamlConfig
)

type YamlConfig struct {
	Console   bool        `yaml:"console"`
	Alert     Alert       `yaml:"alert"`
	Mount     []Mount     `yaml:"mount"`
	Directory []Directory `yaml:"directory"`
	File      []File      `yaml:"file"`
}

type Alert struct {
	Enable bool     `yaml:"enable"`
	Url    string   `yaml:"url"`
	Labels []string `yaml:"labels"`
}

type Mount struct {
	Mount     string   `yaml:"mount"`
	Threshold int      `yaml:"threshold"`
	Scripts   []string `yaml:"scripts"`
	Alert     bool     `yaml:"alert"`
	Alertname string   `yaml:"alertname"`
	Labels    []string `yaml:"labels"`
}

type Directory struct {
	Directory string   `yaml:"directory"`
	Threshold string   `yaml:"threshold"`
	Scripts   []string `yaml:"scripts"`
	Alert     bool     `yaml:"alert"`
	Alertname string   `yaml:"alertname"`
	Labels    []string `yaml:"labels"`
}

type File struct {
	File      string   `yaml:"file"`
	Threshold string   `yaml:"threshold"`
	Scripts   []string `yaml:"scripts"`
	Alert     bool     `yaml:"alert"`
	Alertname string   `yaml:"alertname"`
	Labels    []string `yaml:"labels"`
}
