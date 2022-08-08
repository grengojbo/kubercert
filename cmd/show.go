/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/grengojbo/kubercert/pkg/cert"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show certificates",
	Long:  `Show kubernetes API certificates.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("show called")
		h := cert.HostInfo{
			Host:       Host,
			Port:       Port,
			ExpireDays: ExpireDays,
		}

		if err := h.GetCerts(timeout); err != nil {
			log.Fatalln(err.Error())
		}
		if err := h.ShowCerts(Format); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
