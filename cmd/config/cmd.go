package config

import (
	"fmt"
	"os"
	"path/filepath"

	appconfig "github.com/hduhelp/hdu-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "manage local configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(newInitCmd(), newShowCmd(), newSetCmd())
}

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "write a default config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := configPath(cmd)
			if err != nil {
				return err
			}
			cfg := appconfig.Default()
			cfg.Endpoint = viper.GetString("endpoint")
			cfg.ACID = viper.GetString("acid")
			cfg.Username = viper.GetString("username")
			cfg.CheckIntervalSeconds = viper.GetInt("checkIntervalSeconds")
			cfg.AutoConnect = viper.GetBool("autoConnect")
			cfg.AutoReconnect = viper.GetBool("autoReconnect")
			cfg.LogLevel = viper.GetString("logLevel")

			if err := appconfig.Save(path, cfg, viper.GetString("net.auth.password")); err != nil {
				return err
			}
			fmt.Printf("config written to %s\n", path)
			return nil
		},
	}

	cmd.Flags().StringP("username", "u", "", "username to persist")
	cmd.Flags().StringP("password", "p", "", "password to persist")
	cmd.Flags().String("endpoint", "", "endpoint to persist")
	cmd.Flags().String("acid", "", "ac_id to persist")
	cmd.Flags().Int("interval", 0, "check interval seconds")
	cmd.Flags().Bool("auto-connect", true, "persist auto connect")
	cmd.Flags().Bool("auto-reconnect", true, "persist auto reconnect")

	_ = viper.BindPFlag("username", cmd.Flags().Lookup("username"))
	_ = viper.BindPFlag("net.auth.password", cmd.Flags().Lookup("password"))
	_ = viper.BindPFlag("endpoint", cmd.Flags().Lookup("endpoint"))
	_ = viper.BindPFlag("acid", cmd.Flags().Lookup("acid"))
	_ = viper.BindPFlag("checkIntervalSeconds", cmd.Flags().Lookup("interval"))
	_ = viper.BindPFlag("autoConnect", cmd.Flags().Lookup("auto-connect"))
	_ = viper.BindPFlag("autoReconnect", cmd.Flags().Lookup("auto-reconnect"))

	return cmd
}

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "print the current config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := configPath(cmd)
			if err != nil {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Print(string(data))
			return nil
		},
	}
}

func newSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "update a single config value with --key and --value",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := configPath(cmd)
			if err != nil {
				return err
			}
			key, _ := cmd.Flags().GetString("key")
			value, _ := cmd.Flags().GetString("value")
			if key == "" {
				return fmt.Errorf("key is required")
			}

			v := viper.New()
			v.SetConfigFile(path)
			v.SetConfigType("yaml")
			if err := v.ReadInConfig(); err != nil {
				return err
			}
			v.Set(key, value)
			return v.WriteConfig()
		},
	}
	cmd.Flags().String("key", "", "config key to update")
	cmd.Flags().String("value", "", "new string value")
	return cmd
}

func configPath(cmd *cobra.Command) (string, error) {
	path := cmd.Root().PersistentFlags().Lookup("config").Value.String()
	if path != "" {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hdu-cli.yaml"), nil
}
