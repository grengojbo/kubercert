/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/grengojbo/kubercert/pkg/cert"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// renewCmd represents the renew command
var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "Certificate rotation",
	Long:  `Kubernetes API certificate rotation.`,
	Run: func(cmd *cobra.Command, args []string) {
		h := cert.HostInfo{
			Host:       Host,
			Port:       Port,
			ExpireDays: ExpireDays,
		}

		if err := h.GetCerts(timeout); err != nil {
			log.Fatalln(err.Error())
		}
		if h.IsExpired() {
			if err := h.ReNew(Command); err != nil {
				log.Fatalln(err.Error())
			}
			log.Errorln("Certificate is expired")
		}

		// if err := h.ShowCerts(Format); err != nil {
		// 	log.Fatalln(err.Error())
		// }
		// log.Infof("Date: %s", h.GetExpire(ExpireDays).Format(time.RFC3339))
	},
}

func init() {
	rootCmd.AddCommand(renewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	renewCmd.Flags().StringVarP(&Command, "command", "c", "", "Command to execute")
}
