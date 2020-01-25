package commands

import (
	"github.com/allegro/bigcache"
	"github.com/bowd/quip-exporter/browser"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
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

		repo := repositories.NewFileRepository(
			viper.GetString("repo.basePath"),
		)
		// Todo: Config cache through config.yml
		cache, err := repositories.NewCacheRepository(repo, bigcache.DefaultConfig(10*time.Minute))
		if err != nil {
			logger.Fatal(err)
			return
		}
		go browser.Run(browser.Config{
			Port: viper.GetString("browser.port"),
			Host: viper.GetString("browser.host"),
		}, cache)

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
