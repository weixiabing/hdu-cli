package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	appconfig "github.com/hduhelp/hdu-cli/internal/config"
	platformservice "github.com/hduhelp/hdu-cli/internal/platform/service"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "service",
	Short: "manage background service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(newInstallCmd(), newEnableCmd(), newDisableCmd(), newStatusCmd())
}

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "install the user service definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			binaryPath, err := os.Executable()
			if err != nil {
				return err
			}

			configPath := cmd.Root().PersistentFlags().Lookup("config").Value.String()
			if configPath == "" {
				configPath = filepath.Join(home, ".hdu-cli.yaml")
			}

			installedPath, err := platformservice.InstallUserUnit(home, platformservice.UnitConfig{
				BinaryPath: binaryPath,
				ConfigPath: configPath,
			})
			if err != nil {
				return err
			}

			fmt.Printf("installed service unit at %s\n", installedPath)
			return runSystemctl("--user", "daemon-reload")
		},
	}
}

func newEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable",
		Short: "enable and start the user service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSystemctl("--user", "enable", "--now", "hdu-cli.service")
		},
	}
}

func newDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "disable and stop the user service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSystemctl("--user", "disable", "--now", "hdu-cli.service")
		},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "show the current user service status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSystemctl("--user", "status", "hdu-cli.service")
		},
	}
}

func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hdu-cli.yaml"), nil
}

func InitConfigFile(path string, cfg appconfig.AppConfig, password string) error {
	return appconfig.Save(path, cfg, password)
}

func runSystemctl(args ...string) error {
	cmd := exec.Command("systemctl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
