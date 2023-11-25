package wifi

import (
	"github.com/spf13/cobra"

	"github.com/jodydadescott/shelly-client/sdk/wifi"
)

type callback interface {
	WriteStdout(any) error
	Wifi() (*wifi.Client, error)
}

func NewCmd(callback callback) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "wifi",
		Short: "WiFi Scan / List AP Clients",
	}

	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "returns list of all available networks",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Wifi()
			if err != nil {
				return err
			}

			results, err := client.Scan(cmd.Context())
			if err != nil {
				return err
			}

			callback.WriteStdout(results)
			return nil
		},
	}

	listAPClientsCmd := &cobra.Command{
		Use:   "list-ap-clients",
		Short: "returns list of clients currently connected to the device's access point",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Wifi()
			if err != nil {
				return err
			}

			results, err := client.ListAPClients(cmd.Context())
			if err != nil {
				return err
			}

			callback.WriteStdout(results)
			return nil
		},
	}

	rootCmd.AddCommand(scanCmd, listAPClientsCmd)

	return rootCmd
}
