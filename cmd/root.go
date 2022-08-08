/*
Copyright © 2022 Oleg Dolya

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var version string
var cfgFile string
var Port int
var Host string
var Format string
var Command string
var ExpireDays int
var timeout int
var DryRun bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubercert",
	Short: "Certificates management",
	Long:  `Kubernetes API certificates management.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, commit string, date string, builtBy string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogging)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubercert.yaml)")

	rootCmd.PersistentFlags().StringVarP(&Host, "host", "H", "localhost", "Kubernetes API host")
	rootCmd.PersistentFlags().IntVarP(&Port, "port", "p", 6443, "Kubernetes API port")
	rootCmd.PersistentFlags().StringVarP(&Format, "output", "o", "text", "Output format (text, json)")
	rootCmd.PersistentFlags().IntVarP(&ExpireDays, "expire", "e", 7, "Certificate expiration days")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 5, "Timeout in seconds")
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "", false, "Dry run")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output (debug logging)")
	rootCmd.PersistentFlags().Bool("trace", false, "Enable super verbose output (trace logging)")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("trace", rootCmd.PersistentFlags().Lookup("trace"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kubercert" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kubercert")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// initLogging initializes the logger
func initLogging() {
	if viper.GetBool("trace") {
		log.SetLevel(log.TraceLevel)
	} else if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	} else {
		switch logLevel := strings.ToUpper(os.Getenv("LOG_LEVEL")); logLevel {
		case "TRACE":
			log.SetLevel(log.TraceLevel)
		case "DEBUG":
			log.SetLevel(log.DebugLevel)
		case "WARN":
			log.SetLevel(log.WarnLevel)
		case "ERROR":
			log.SetLevel(log.ErrorLevel)
		default:
			log.SetLevel(log.InfoLevel)
		}
	}
	log.SetOutput(ioutil.Discard)
	log.AddHook(&writer.Hook{
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{
		Writer: os.Stdout,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
			log.TraceLevel,
		},
	})

	formatter := &log.TextFormatter{
		ForceColors: true,
	}

	log.SetFormatter(formatter)

}
