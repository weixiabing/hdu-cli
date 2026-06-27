package net

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/hduhelp/hdu-cli/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newDaemonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "run reconnect loop in foreground",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bindCredentialFlags(cmd); err != nil {
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

	cmd.Flags().StringP("username", "u", "", "username of srun")
	cmd.Flags().StringP("password", "p", "", "password of srun")
	cmd.Flags().IntP("interval", "i", 60, "second interval of daemon mode")

	return cmd
}

func runDaemon(ctx context.Context) error {
	if sessionService == nil {
		return errors.New("network service is not initialized")
	}

	creds, err := currentCredentials()
	if err != nil {
		return err
	}

	checkInterval := time.Duration(viper.GetInt("checkIntervalSeconds")) * time.Second
	if checkInterval <= 0 {
		checkInterval = time.Minute
	}

	manager := core.NewReconnectManager(sessionService, core.ReconnectConfig{
		Interval:   time.Second,
		MaxBackoff: checkInterval,
	})
	manager.OnStateChange(func(status core.Status) {
		log.Printf("daemon state=%s online=%t message=%q", status.Phase, status.Online, status.Message)
	})

	log.Printf("start daemon: check every %d seconds", int(checkInterval/time.Second))
	if err := ensureConnected(ctx, manager, creds); err != nil {
		return err
	}

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := ensureConnected(ctx, manager, creds); err != nil {
				return err
			}
			if portalServer.InternetReachable(ctx) {
				continue
			}

			log.Println("internet is not available, reconnecting")
			if err := sessionService.Logout(ctx, creds); err != nil {
				log.Printf("logout before reconnect failed: %v", err)
			}
			if err := ensureConnected(ctx, manager, creds); err != nil {
				return err
			}
		}
	}
}

func ensureConnected(ctx context.Context, manager *core.ReconnectManager, creds core.Credentials) error {
	status, err := sessionService.CurrentStatus(ctx)
	if err == nil && status.Online {
		return nil
	}

	status, err = manager.ReconnectOnce(ctx, creds)
	if err != nil {
		return err
	}
	if !status.Online {
		return errors.New("reconnect finished without an online session")
	}
	return nil
}
