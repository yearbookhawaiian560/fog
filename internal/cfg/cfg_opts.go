package cfg

import (
	"github.com/spf13/pflag"
)

type CfgOption func(*ConfigManager)

func WithFlags(flagSet *pflag.FlagSet) CfgOption {
	return func(mgr *ConfigManager) {
		mgr.flags = flagSet
	}
}
