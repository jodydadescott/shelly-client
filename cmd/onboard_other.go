//go:build !linux

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (t *Cmd) newOnboardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "onboard",
		Short: "Linux only: joins Shelly access points to set auth and wifi",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("onboard needs Linux (nl80211)")
		},
	}
}
