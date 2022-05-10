package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes the configuration",
	Long:  `This subcommand generates config file with all path to MSP related files`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return
		}
		viper.Set("conn_config", fmt.Sprintf("%s/organizations/peerOrganizations/org1.example.com/connection-org1.yaml", path))
		viper.Set("cred_path", fmt.Sprintf("%s/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp", path))
		err = viper.WriteConfig()
		if err != nil {
			return
		}
		fmt.Println("Successfully initialized config")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("path", "p", "", "Path to the test Fabric test network directory")
	err := initCmd.MarkFlagRequired("path")
	if err != nil {
		return
	}
}
