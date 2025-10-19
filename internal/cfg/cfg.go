package cfg

import (
	"log/slog"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type ConfigManager struct {
	// Config Providers, all implement koanf.Provider interface
	envProvider koanf.Provider

	// Configurable attributes
	flags *pflag.FlagSet
}

// Globals to be initialized
var (
	// Main config to access vars
	K *koanf.Koanf
	// Config manager to load, watch and potentially reload the main config
	Mgr *ConfigManager
)

func Init(opts ...CfgOption) {
	if Mgr != nil {
		slog.Warn("Attempting to init config again")
	}

	K = koanf.New(".")

	Mgr = &ConfigManager{
		envProvider: env.Provider("", ".", func(s string) string {
			return strings.Replace(s, "__", ".", -1)
		}),
	}

	for _, opt := range opts {
		opt(Mgr)
	}

	if err := load(); err != nil {
		slog.Error("Error loading config", "err", err)
		panic(err)
	}

	slog.Debug("All cfg", "cfg", K.All())
}

func load() error {
	if Mgr.envProvider != nil {
		if err := K.Load(Mgr.envProvider, nil); err != nil {
			return err
		}
	}

	if Mgr.flags != nil {
		if err := K.Load(posflag.Provider(Mgr.flags, ".", K), nil); err != nil {
			return err
		}
	}

	return nil
}
