package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/0div/fog/cmd/migrate"
	"github.com/0div/fog/cmd/repl"
	"github.com/0div/fog/cmd/seed"
	"github.com/0div/fog/internal/cfg"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:              "fog",
		Short:            "fog",
		Long:             "fog",
		PersistentPreRun: bindFlagsAndInitConfig,
	}

	rootCmd.PersistentFlags().BoolP("development", "d", false, "Development mode (prints prettier log messages)")
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Debug mode (prints debug messages and call traces)")

	replCmd := repl.Setup()
	rootCmd.AddCommand(replCmd)
	rootCmd.AddCommand(migrate.Setup())
	rootCmd.AddCommand(seed.Setup())

	cmd, _, err := rootCmd.Find(os.Args[1:])
	// default to repl if no cmd is given
	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{replCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func bindFlagsAndInitConfig(cmd *cobra.Command, args []string) {
	viper.BindPFlag("development", cmd.Flags().Lookup("development"))
	viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	if env == "local" {
		godotenv.Load(".env.local")
	} else {
		godotenv.Load(".env." + env)
	}
	godotenv.Load(".env." + env)

	cfg.Init(cfg.WithFlags(cmd.Flags()))

	cfg.K.Set("env", env)

	// Setup default global logger
	var logLevel = slog.LevelInfo
	var addSource = false

	if cfg.Bool("debug") {
		fmt.Println("debug mode enabled")
		logLevel = slog.LevelDebug
		addSource = true
	}

	var logHandler slog.Handler
	if env == "local" {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: addSource,
		})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: addSource,
		})
	}

	slog.SetDefault(slog.New(logHandler))
}
