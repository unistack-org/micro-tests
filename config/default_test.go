package config

import (
	"os"
	"testing"

	envconfig "github.com/unistack-org/micro-config-env/v3"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/config"
)

type Config struct {
	String string `env:"MICRO_TEST" default:"default"`
}

func TestMultiple(t *testing.T) {
	cfg := &Config{}
	svc := micro.NewService(micro.Configs(
		config.NewConfig(config.Struct(cfg)),
		envconfig.NewConfig(config.Struct(cfg)),
	),
	)
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if cfg.String != "default" {
		t.Fatalf("config not parsed by default source: %#+v\n", cfg)
	}
	os.Setenv("MICRO_TEST", "non_default")
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if cfg.String == "default" {
		t.Fatalf("config not parsed by default source: %#+v\n", cfg)
	}
	if cfg.String != "non_default" {
		t.Fatalf("config not parsed by default source: %#+v\n", cfg)
	}
	t.Logf("config: %#+v\n", cfg)
}
