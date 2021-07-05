package config

import (
	"context"
	"os"
	"testing"

	envconfig "github.com/unistack-org/micro-config-env/v3"
	"github.com/unistack-org/micro/v3/config"
)

type Config struct {
	String string `env:"MICRO_TEST" default:"default"`
}

func TestMultiple(t *testing.T) {
	ctx := context.Background()
	cfg := &Config{}

	c1 := config.NewConfig(config.Struct(cfg))
	c2 := envconfig.NewConfig(config.Struct(cfg))

	if err := c1.Init(); err != nil {
		t.Fatal(err)
	}
	if err := c2.Init(); err != nil {
		t.Fatal(err)
	}

	if err := c1.Load(ctx); err != nil {
		t.Fatal(err)
	}
	if err := c2.Load(ctx); err != nil {
		t.Fatal(err)
	}

	if cfg.String != "default" {
		t.Fatalf("config not parsed by default source: %#+v\n", cfg)
	}
	os.Setenv("MICRO_TEST", "non_default")
	if err := c1.Load(ctx, config.LoadOverride(true)); err != nil {
		t.Fatal(err)
	}
	if err := c2.Load(ctx, config.LoadOverride(true)); err != nil {
		t.Fatal(err)
	}
	if cfg.String == "default" {
		t.Fatalf("config not parsed by default source: %#+v\n", cfg)
	}
	if cfg.String != "non_default" {
		t.Fatalf("config not parsed by default source: %#+v\n", cfg)
	}
}
