package config

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"flag"
)

const (
	ctxKey configCtxKey = "ctx_key_for_config"
)

var (
	configPath = flag.String("config", "", "path to config file")
)

type configCtxKey string

// Config is a struct for config
type Config struct {
	DbPath    string `yaml:"db.path"`
	DataPath  string `yaml:"data.path"`
	ListenUrl string `yaml:"listen.url"`
}

// NewContext places config into context
func NewContext(ctx context.Context, config interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := config.(*Config); !ok {
		flag.Parse()
		config = NewConfig(*configPath)
	}

	return context.WithValue(ctx, ctxKey, config)
}

// FromContext returns config from context
func FromContext(ctx context.Context) *Config {
	if config, ok := ctx.Value(ctxKey).(*Config); ok {
		return config
	}

	flag.Parse()
	return NewConfig(*configPath)
}

// NewConfig loads config from file
func NewConfig(filePath string) *Config {
	file, err := os.Open(filePath)
	if err != nil {
		log.Panicf("erorr when open config.yaml: %s", err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicf("erorr when reading config: %s", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Panicf("erorr when unmarshall congig: %s", err)
	}

	return config
}
