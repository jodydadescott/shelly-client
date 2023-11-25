package switchx

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jodydadescott/shelly-client/sdk/switchx"
)

var (
	truePointer  = true
	falsePointer = false
)

type callback interface {
	Switch() (*switchx.Client, error)
}

func NewCmd(callback callback) *cobra.Command {

	var switchIDArg string

	getSwitchID := func() (*int, error) {

		if switchIDArg == "" {
			return nil, fmt.Errorf("switchID is required")
		}

		switchID, err := strconv.Atoi(switchIDArg)
		if err == nil {
			return &switchID, nil
		}

		return nil, fmt.Errorf("switchID must be an integer")

	}

	rootCmd := &cobra.Command{
		Use:   "switch",
		Short: "Turn switch on or off",
	}

	rootCmd.PersistentFlags().StringVar(&switchIDArg, "id", "", "switch ID integer")

	setOnCmd := &cobra.Command{
		Use:   "on",
		Short: "Turn light on",
		RunE: func(cmd *cobra.Command, args []string) error {

			switchID, err := getSwitchID()
			if err != nil {
				return err
			}

			client, err := callback.Switch()
			if err != nil {
				return err
			}

			return client.Set(cmd.Context(), *switchID, &truePointer)
		},
	}

	setOffCmd := &cobra.Command{
		Use:   "off",
		Short: "Turn light off",
		RunE: func(cmd *cobra.Command, args []string) error {

			switchID, err := getSwitchID()
			if err != nil {
				return err
			}

			client, err := callback.Switch()
			if err != nil {
				return err
			}

			return client.Set(cmd.Context(), *switchID, &falsePointer)
		},
	}

	toggleCmd := &cobra.Command{
		Use:   "toggle",
		Short: "Toggles switch",
		RunE: func(cmd *cobra.Command, args []string) error {

			switchID, err := getSwitchID()
			if err != nil {
				return err
			}

			client, err := callback.Switch()
			if err != nil {
				return err
			}

			return client.Toggle(cmd.Context(), *switchID)
		},
	}

	rootCmd.AddCommand(setOnCmd, setOffCmd, toggleCmd)
	return rootCmd
}
