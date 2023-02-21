package config

type CacheConfig struct {
	Type    string `json:"type" yaml:"type"`
	Options map[string]ConfigList
}
