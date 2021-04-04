package cmd

import (
	"fmt"
	"os"

	"github.com/lbryio/notifica/app/store"

	"github.com/lbryio/notifica/app/action"

	"github.com/lbryio/notifica/app/config"
	"github.com/lbryio/notifica/app/env"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	//cobra.OnInitialize(config.InitializeConfiguration)
	rootCmd.PersistentFlags().BoolVarP(&config.IsDebugMode, "debugmode", "d", true, "turns on debug mode for the application command.")
	rootCmd.PersistentFlags().BoolP("tracemode", "t", false, "turns on trace mode for the application command, very verbose logging.")
	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		logrus.Panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "notifica",
	Short: "Notifications system for LBRY Inc",
	Long:  `Creates/Stores/Returns Device, Web, Email, and App Notifications`,
	Run: func(cmd *cobra.Command, args []string) {
		//Run the application
		conf, err := env.NewWithEnvVars()
		if err != nil {
			logrus.Panic(err)
		}
		config.InitLogging(conf)
		store.Init(conf)
		action.StartNotifica(conf.APIServerPort)
	},
}

// Execute executes the root command and is the entry point of the application from main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
