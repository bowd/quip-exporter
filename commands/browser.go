package commands

import (
	"github.com/bowd/quip-exporter/browser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var browserCmd = &cobra.Command{
	Use:   "browser",
	Short: "start browser",
	Long:  "Start the Quip archive browser",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.WithField("module", "main")
		stopChan := make(chan os.Signal, 1)
		doneChan := make(chan bool, 1)
		signal.Notify(stopChan, syscall.SIGINT)
		signal.Notify(stopChan, syscall.SIGTERM)

		go browser.Run(browser.Config{
			Port: viper.GetString("browser.port"),
			Host: viper.GetString("browser.host"),
		})

		cleanup := func() {
			// Cleanup here
		}

		select {
		case <-stopChan:
			logger.Info("Got stop signal. Finishing in-flight work.")
			cleanup()
			logger.Info("Work done. Goodbye!")
		case <-doneChan:
			cleanup()
			logger.Info("Done scraping. Goodbye!")
		}
	},
}

func init() {
}
