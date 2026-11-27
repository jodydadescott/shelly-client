package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mdlayher/wifi"
	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"

	"github.com/jodydadescott/shelly-client/cmd/util"
	sdk_client "github.com/jodydadescott/shelly-client/sdk/client"
	"github.com/jodydadescott/shelly-client/sdk/devconfig"
)

const (
	// A Shelly device in access point mode is 192.168.33.1/24.
	shellyAPDevice  = "192.168.33.1"
	shellyAPAddress = "192.168.33.200/24"
)

func (t *Cmd) newOnboardCmd() *cobra.Command {

	var (
		ifname   string
		region   string
		interval time.Duration
		cooldown time.Duration
	)

	cmd := &cobra.Command{
		Use:   "onboard",
		Short: "Runs until interrupted: joins the access point of any Shelly device that has a device-ID config, sets its auth and wifi, then leaves",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}
			if config.Shelly == nil {
				return fmt.Errorf("shelly config is required")
			}

			var ids []string
			for id := range config.Shelly.ShellyConfigs {
				ids = append(ids, id)
			}

			wc, err := wifi.New()
			if err != nil {
				return err
			}
			defer wc.Close()

			if region != "" {
				if err := wc.SetRegulatoryRegion(region, wifi.RegulatoryHintUser); err != nil {
					return fmt.Errorf("set regulatory region %s: %w", region, err)
				}
			}

			ifi, err := wifiInterface(wc, ifname)
			if err != nil {
				return err
			}
			link, err := netlink.LinkByName(ifi.Name)
			if err != nil {
				return err
			}
			if err := netlink.LinkSetUp(link); err != nil {
				return err
			}

			nextTry := map[string]time.Time{}

			pass := func() error {
				scanCtx, scanCancel := context.WithTimeout(ctx, 20*time.Second)
				defer scanCancel()
				if err := wc.Scan(scanCtx, ifi); err != nil {
					return fmt.Errorf("scan: %w", err)
				}
				aps, err := wc.AccessPoints(ifi)
				if err != nil {
					return err
				}

				for _, ap := range aps {
					id, ok := util.DeviceIDForSSID(ap.SSID, ids)
					if !ok || time.Now().Before(nextTry[id]) {
						continue
					}
					nextTry[id] = time.Now().Add(cooldown)

					log := zap.L().With(zap.String("ssid", ap.SSID), zap.String("deviceID", id))
					if err := onboardDevice(ctx, wc, ifi, link, ap.SSID, id, config); err != nil {
						log.Error("onboard failed", zap.Error(err))
						continue
					}
					log.Warn("onboarded")
				}
				return nil
			}

			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			for {
				if err := pass(); err != nil {
					zap.L().Error("onboard pass failed", zap.Error(err))
				}
				select {
				case <-ctx.Done():
					return nil
				case <-ticker.C:
				}
			}
		},
	}

	cmd.Flags().StringVar(&ifname, "interface", "", "Wi-Fi interface to scan and join with (empty: the only one there is)")
	cmd.Flags().StringVar(&region, "region", "", "regulatory region to set first, for example US (empty leaves it unchanged)")
	cmd.Flags().DurationVar(&interval, "interval", 30*time.Second, "time between scans")
	cmd.Flags().DurationVar(&cooldown, "cooldown", 10*time.Minute, "minimum time between attempts on the same device")
	return cmd
}

// wifiInterface returns the Wi-Fi interface called name, or with name "", the only one there is.
func wifiInterface(wc *wifi.Client, name string) (*wifi.Interface, error) {
	ifis, err := wc.Interfaces()
	if err != nil {
		return nil, err
	}
	var found []*wifi.Interface
	for _, ifi := range ifis {
		if ifi.Name != "" && (name == "" || ifi.Name == name) {
			found = append(found, ifi)
		}
	}
	if len(found) != 1 {
		return nil, fmt.Errorf("want exactly one Wi-Fi interface named %q (empty: any), found %d", name, len(found))
	}
	return found[0], nil
}

// onboardDevice joins the device's open access point, checks it is device id, sets the desired
// auth and then the desired wifi (last, since it may turn the access point off), and leaves.
// The rest of the desired config is left to `reconcile`, once the device is on the LAN.
func onboardDevice(ctx context.Context, wc *wifi.Client, ifi *wifi.Interface, link netlink.Link, ssid, id string, config *Config) error {

	if err := wc.Connect(ifi, ssid); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer wc.Disconnect(ifi)

	if err := waitAssociated(ctx, wc, ifi, 15*time.Second); err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(shellyAPAddress)
	if err != nil {
		return err
	}
	if err := netlink.AddrReplace(link, addr); err != nil {
		return err
	}
	defer netlink.AddrDel(link, addr)

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	client := util.NewShellyClient(config, shellyAPDevice)
	defer client.Close()

	info, err := client.GetDeviceInfo(ctx)
	if err != nil {
		return err
	}
	if !strings.EqualFold(*info.ID, id) {
		return fmt.Errorf("access point %s is device %s, not %s", ssid, *info.ID, id)
	}

	desired, err := sdk_client.GetConfig(ctx, client)
	if err != nil {
		return err
	}
	wifiConfig, ok := desired["wifi"]
	if !ok {
		return fmt.Errorf("desired config has no wifi")
	}

	if auth, ok := desired[devconfig.AuthKey]; ok {
		if err := client.SetAuth(ctx, auth); err != nil {
			return fmt.Errorf("set auth: %w", err)
		}
	}

	methods, err := client.ListMethods(ctx)
	if err != nil {
		return err
	}
	method, params, err := devconfig.SetConfigCall("wifi", wifiConfig, methods.Methods)
	if err != nil {
		return err
	}
	if _, err := client.Call(ctx, method, params); err != nil {
		return fmt.Errorf("set wifi: %w", err)
	}
	return nil
}

func waitAssociated(ctx context.Context, wc *wifi.Client, ifi *wifi.Interface, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if bss, err := wc.BSS(ifi); err == nil && bss.Status == wifi.BSSStatusAssociated {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
	return fmt.Errorf("not associated after %s", timeout)
}
