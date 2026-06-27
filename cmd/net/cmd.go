package net

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hduhelp/hdu-cli/internal/core"
	"github.com/hduhelp/hdu-cli/pkg/srun"
	"github.com/hduhelp/hdu-cli/pkg/table"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cmd represents the srun command
var Cmd = &cobra.Command{
	Use:   "net",
	Short: "i-hdu network auth cli",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := url.ParseRequestURI(viper.GetString("net.endpoint")); err != nil {
			return err
		}
		portalServer = srun.New(viper.GetString("net.endpoint"), viper.GetString("net.acid"))
		sessionService = core.NewSessionService(portalServer)
		return nil
	},
}

var portalServer *srun.PortalServer
var sessionService *core.SessionService

func init() {
	Cmd.PersistentFlags().StringP("endpoint", "e", "", "endpoint host of srun")
	viper.SetDefault("net.endpoint", "http://192.168.112.30")
	cobra.CheckErr(viper.BindPFlag("net.endpoint", Cmd.PersistentFlags().Lookup("endpoint")))

	Cmd.PersistentFlags().StringP("acid", "a", "", "ac_id of srun")
	viper.SetDefault("net.acid", detectDefaultACID())
	cobra.CheckErr(viper.BindPFlag("net.acid", Cmd.PersistentFlags().Lookup("acid")))

	loginCmd.Flags().StringP("username", "u", "", "username of srun")
	loginCmd.Flags().StringP("password", "p", "", "password of srun")
	loginCmd.Flags().BoolP("daemon", "d", false, "daemon mode")
	loginCmd.Flags().IntP("interval", "i", 60, "second interval of daemon mode")

	logoutCmd.Flags().StringP("username", "u", "", "username of srun")
	cobra.CheckErr(viper.BindPFlag("net.auth.username", logoutCmd.Flags().Lookup("username")))

	Cmd.AddCommand(infoCmd, statusCmd, loginCmd, logoutCmd, internetCmd, newDaemonCmd())
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "show info of your i-hdu network",
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := portalServer.GetUserInfo()
		if err != nil {
			return err
		}
		table.PrintStruct(info, "chinese")
		return nil
	},
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login i-hdu of the account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := bindCredentialFlags(cmd); err != nil {
			return err
		}
		creds, err := currentCredentials()
		if err != nil {
			return err
		}
		status, err := sessionService.Login(cmd.Context(), creds)
		if err != nil {
			return err
		}
		printStatus(status)

		runAsDaemon, err := cmd.Flags().GetBool("daemon")
		if err != nil || !runAsDaemon {
			return err
		}

		interval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			return err
		}
		if interval > 0 {
			viper.Set("checkIntervalSeconds", interval)
		}
		return runDaemon(cmd.Context())
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show the current network session state",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, err := sessionService.CurrentStatus(cmd.Context())
		if err != nil && status.Phase == "" {
			return err
		}
		printStatus(status)
		return err
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout i-hdu of the account",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, err := sessionService.CurrentStatus(cmd.Context())
		if err == nil && !status.Online {
			fmt.Println("you are not login")
			return nil
		}
		if err != nil && !errors.Is(err, core.ErrAuthenticationFailed) {
			return err
		}

		username, err := currentUsername(context.Background())
		if err != nil {
			return err
		}
		if viper.GetBool("verbose") {
			info, infoErr := portalServer.GetUserInfo()
			if infoErr == nil {
				table.PrintStruct(info, "chinese")
			}
		}
		fmt.Printf("you are logout account %s\n", username)
		if err := sessionService.Logout(cmd.Context(), core.Credentials{Username: username}); err != nil {
			return err
		}
		return nil
	},
}

// internetCmd represents the logout command
var internetCmd = &cobra.Command{
	Use:   "internet",
	Short: "check if connect to the internet",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(portalServer.InternetReachable(cmd.Context()))
		return nil
	},
}

func detectDefaultACID() string {
	resp, _, errs := gorequest.New().Get("http://www.baidu.com").End()
	if errs != nil || resp == nil || resp.Request == nil || resp.Request.URL == nil {
		return "0"
	}
	if acid := resp.Request.URL.Query().Get("ac_id"); acid != "" {
		return acid
	}
	return "0"
}

func bindCredentialFlags(cmd *cobra.Command) error {
	if err := viper.BindPFlag("net.auth.username", cmd.Flags().Lookup("username")); err != nil {
		return err
	}
	if err := viper.BindPFlag("net.auth.password", cmd.Flags().Lookup("password")); err != nil {
		return err
	}
	return nil
}

func currentCredentials() (core.Credentials, error) {
	creds := core.Credentials{
		Username: viper.GetString("net.auth.username"),
		Password: viper.GetString("net.auth.password"),
	}
	if creds.Username == "" {
		return core.Credentials{}, errors.New("username is empty")
	}
	if creds.Password == "" {
		return core.Credentials{}, errors.New("password is empty")
	}
	return creds, nil
}

func currentUsername(ctx context.Context) (string, error) {
	if username := viper.GetString("net.auth.username"); username != "" {
		return username, nil
	}
	info, err := portalServer.GetUserInfo()
	if err != nil {
		return "", err
	}
	if info.UserName == "" {
		return "", errors.New("username is empty")
	}
	return info.UserName, nil
}

func printStatus(status core.Status) {
	if status.Message != "" {
		fmt.Printf("%s: %s\n", status.Phase, status.Message)
		return
	}
	fmt.Println(status.Phase)
}
