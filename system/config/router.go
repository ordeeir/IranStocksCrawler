package config

type RouterConfig struct {
	Type    string `json:"type" yaml:"type"`
	Options map[string]ConfigList
}
