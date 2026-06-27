package cmd

import (
	"fmt"
	"os"

	"github.com/hduhelp/hdu-cli/cmd/auth"
	configcmd "github.com/hduhelp/hdu-cli/cmd/config"
	"github.com/hduhelp/hdu-cli/cmd/net"
	"github.com/hduhelp/hdu-cli/cmd/rpc"
	servicecmd "github.com/hduhelp/hdu-cli/cmd/service"
	appconfig "github.com/hduhelp/hdu-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "hdu_cli",
	Short:   "hdu cli",
	Version: "alpha",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("save") {
			if viper.WriteConfig() != nil {
				viper.Set("verbose", nil)
				viper.Set("save", nil)
				cobra.CheckErr(viper.SafeWriteConfig())
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hdu-cli.yaml)")
	rootCmd.PersistentFlags().BoolP("save", "s", false, "save config")
	cobra.CheckErr(viper.BindPFlag("save", rootCmd.PersistentFlags().Lookup("save")))
	rootCmd.PersistentFlags().BoolP("verbose", "V", false, "show more info")
	cobra.CheckErr(viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")))

	rootCmd.AddCommand(net.Cmd, auth.Cmd, rpc.Cmd, configcmd.Cmd, servicecmd.Cmd)
}

var cfgFile string

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	defaults := appconfig.Default()
	viper.RegisterAlias("endpoint", "net.endpoint")
	viper.RegisterAlias("acid", "net.acid")
	viper.RegisterAlias("username", "net.auth.username")
	viper.SetDefault("endpoint", defaults.Endpoint)
	viper.SetDefault("acid", defaults.ACID)
	viper.SetDefault("username", defaults.Username)
	viper.SetDefault("checkIntervalSeconds", defaults.CheckIntervalSeconds)
	viper.SetDefault("autoConnect", defaults.AutoConnect)
	viper.SetDefault("autoReconnect", defaults.AutoReconnect)
	viper.SetDefault("logLevel", defaults.LogLevel)
	viper.SetDefault("net.endpoint", defaults.Endpoint)
	viper.SetDefault("net.acid", defaults.ACID)
	viper.SetDefault("net.auth.username", defaults.Username)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".hdu_cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hdu-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
