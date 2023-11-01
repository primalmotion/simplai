package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfgName string
)

var (
	version = "v0.0.0"
	commit  = "dev"
)

func main() {

	cobra.OnInitialize(initCobra)
	mainCtx := context.Background()

	rootCmd := &cobra.Command{
		Use:              "simplai",
		Short:            "fairely usable AI assistant based on simplai",
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
				return err
			}
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetBool("version") {
				fmt.Printf("simplai %s (%s)\n", version, commit)
				os.Exit(0)
				return nil
			}

			engine := viper.GetString("engine")
			model := viper.GetString("model")
			api := viper.GetString("api")
			searxURL := viper.GetString("searxurl")
			debug := viper.GetBool("debug")

			return run(cmd.Context(), engine, model, api, searxURL, debug)
		},
	}
	rootCmd.Flags().Bool("version", false, "Show version")
	rootCmd.Flags().String("engine", "openai", "Select the engine to use (openai or ollama)")
	rootCmd.Flags().String("api", "", "Set the server API base url")
	rootCmd.Flags().String("model", "HuggingFaceH4/zephyr-7b-beta", "Select the model to use")
	rootCmd.Flags().String("searx-url", "", "Set the searx url")
	rootCmd.Flags().Bool("debug", false, "Set the default debug mode state")

	if err := rootCmd.ExecuteContext(mainCtx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func initCobra() {

	viper.SetEnvPrefix("simplai")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln("unable to find home dir: ", err)
	}

	if cfgFile == "" {
		cfgFile = os.Getenv("SIMPLAI_CONFIG")
	}

	if cfgFile != "" {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			log.Fatalln("config file does not exist", err)
		}

		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfgFile)

		if err = viper.ReadInConfig(); err != nil {
			log.Fatalln("unable to read config", cfgFile)
		}

		return
	}

	viper.AddConfigPath(path.Join(home, ".config", "simplai"))
	viper.AddConfigPath("/usr/local/etc/simplai")
	viper.AddConfigPath("/etc/simplai")

	if cfgName == "" {
		cfgName = os.Getenv("SIMPLAI_CONFIG_NAME")
	}

	if cfgName == "" {
		cfgName = "config"
	}

	viper.SetConfigName(cfgName)

	if err = viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Fatalln("unable to read config:", err)
		}
	}
}
