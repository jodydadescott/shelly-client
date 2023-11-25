package light

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jodydadescott/shelly-client/sdk/light"
)

var (
	truePointer  = true
	falsePointer = false
)

type callback interface {
	Light() (*light.Client, error)
	GetCTX() context.Context
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
		Use:   "light",
		Short: "Turn light on, off, or set brigtness level",
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

			client, err := callback.Light()
			if err != nil {
				return err
			}

			return client.Set(callback.GetCTX(), *switchID, &truePointer, nil)
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

			client, err := callback.Light()
			if err != nil {
				return err
			}

			return client.Set(callback.GetCTX(), *switchID, &falsePointer, nil)
		},
	}

	setBrightnessCmd := &cobra.Command{
		Use:   "bright",
		Short: "Sets light brightness",
		RunE: func(cmd *cobra.Command, args []string) error {

			switchID, err := getSwitchID()
			if err != nil {
				return err
			}

			client, err := callback.Light()
			if err != nil {
				return err
			}

			brightness := 0.0

			if len(args) > 0 {
				arg := args[0]
				f, err := strconv.ParseFloat(arg, 32)
				if err != nil {
					return fmt.Errorf("arg %s is not a valid float", arg)
				}
				brightness = f
			}

			return client.Set(callback.GetCTX(), *switchID, nil, &brightness)
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

			client, err := callback.Light()
			if err != nil {
				return err
			}

			return client.Toggle(callback.GetCTX(), *switchID)
		},
	}

	rootCmd.AddCommand(toggleCmd, setOnCmd, setOffCmd, setBrightnessCmd)
	return rootCmd
}
