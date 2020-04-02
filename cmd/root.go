package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tianhongw/misc-go/conf"
	"github.com/tianhongw/misc-go/db"
	"github.com/tianhongw/misc-go/log"
)

const (
	defaultCfgFile = "$HOME/.app.toml"
	defaultCfgType = "toml"
)

var (
	cfgFile string
	cfgType string
)

var rootCmd = cobra.Command{
	Use: "app",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("config file used: %s", cfgFile)
	},
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	cobra.OnInitialize(doInit)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", fmt.Sprintf("config file path, default use: %s", defaultCfgFile))
	rootCmd.PersistentFlags().StringVarP(&cfgType, "type", "t", "", fmt.Sprintf("config file type, default use: %s", defaultCfgType))
	rootCmd.PersistentFlags().StringP("log", "l", "", "log file path")

	viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))
}

func doInit() {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get home dir: %v", err)
			os.Exit(1)
		}
		cfgFile = strings.Replace(defaultCfgFile, "$HOME", home, 1)
	}

	if cfgType == "" {
		cfgType = defaultCfgType
	}

	if cfgFileUsed, err := conf.Init(cfgFile, cfgType); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config file: %v", err)
		os.Exit(1)
	} else {
		cfgFile = cfgFileUsed
	}
}

func run() {
	if err := conf.Opts.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %v", err)
		os.Exit(1)
	}

	if err := log.Init(conf.Opts); err != nil {
		fmt.Fprintf(os.Stderr, "init log failed: %v", err)
		os.Exit(1)
	}
	defer log.Flush()

	if err := db.Init(conf.Opts); err != nil {
		log.Errorf("init database failed: %v", err)
		os.Exit(1)
	}
	defer db.Close()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "execute failed: %v", err)
		os.Exit(1)
	}
}
