package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/jodydadescott/shelly-client/cmd/util"
	sdk_client "github.com/jodydadescott/shelly-client/sdk/client"
)

func (t *Cmd) newReconcileCmd() *cobra.Command {

	var interval time.Duration

	cmd := &cobra.Command{
		Use:   "reconcile",
		Short: "Runs until interrupted: every interval, sets each device's desired config where it has drifted, and logs the drift",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyStatus) error {

				desired, err := sdk_client.GetConfig(ctx, client)
				if err != nil {
					return err
				}

				report, err := client.ApplyDeviceConfig(ctx, desired, false)
				if err != nil {
					return err
				}

				if report.NoChange {
					zap.L().Debug("in sync", zap.String("hostname", hostname), zap.String("deviceID", *deviceInfo.ID))
					return nil
				}

				// Uptime is from before the change: a small value next to a large drift suggests the
				// device reset itself, as opposed to someone changing one setting.
				var uptime float64
				if deviceStatus != nil && deviceStatus.System != nil && deviceStatus.System.Uptime != nil {
					uptime = *deviceStatus.System.Uptime
				}

				zap.L().Warn("drift corrected",
					zap.String("hostname", hostname),
					zap.String("deviceID", *deviceInfo.ID),
					zap.String("app", *deviceInfo.App),
					zap.Float64("uptimeSeconds", uptime),
					zap.Strings("drift", report.Drift),
					zap.Bool("rebooted", report.RebootRequired))
				return nil
			}

			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			for {
				if err := util.Process(ctx, config, "reconcile", true, do); err != nil {
					zap.L().Error("reconcile pass had errors", zap.Error(err))
				}

				select {
				case <-ctx.Done():
					return nil
				case <-ticker.C:
				}
			}
		},
	}

	cmd.Flags().DurationVar(&interval, "interval", 5*time.Minute, "time between passes")
	return cmd
}
