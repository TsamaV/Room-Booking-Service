package configs_test

import (
	"Avito/configs"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DSN", "postgres://test")
	os.Setenv("SECRET", "secret")

	cfg := configs.LoadConfig()

	if cfg.Db.Dsn != "postgres://test" {
		t.Fatalf("expected postgres://test, got %s", cfg.Db.Dsn)
	}
	if cfg.Auth.Secret != "secret" {
		t.Fatalf("expected secret, got %s", cfg.Auth.Secret)
	}
}