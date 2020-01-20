package commands

import (
	"context"
	"github.com/bowd/quip-exporter/scraper"
	"github.com/bowd/quip-exporter/utils"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bowd/quip-exporter/client"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Say hello!",
	Long:  "Address a wonderful greeting to the majestic executioner of this CLI",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.WithField("module", "main")
		err := utils.EnsureDir("./output")
		if err != nil {
			logger.Errorln(err)
			return
		}

		stopChan := make(chan os.Signal, 1)
		doneChan := make(chan bool, 1)
		signal.Notify(stopChan, syscall.SIGINT)
		signal.Notify(stopChan, syscall.SIGTERM)
		ctx, cancel := context.WithCancel(context.Background())
		db, err := bolt.Open(viper.GetString("scraper.dbpath"), 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			logger.Errorf("Could not open database: %s", err)
		}
		quipClient, err := client.New(
			viper.GetString("scraper.token"),
			viper.GetInt("scraper.tokenConcurrency"),
			viper.GetInt("scraper.rps"),
			viper.GetDuration("scraper.batch.wait"),
			viper.GetInt("scraper.batch.maxItems"),
		)
		if err != nil {
			logger.Errorln(err)
			return
		}

		scraper := scraper.New(
			quipClient,
			db,
			viper.GetStringSlice("scraper.folders"),
		)

		go scraper.Run(ctx, doneChan)

		cleanup := func() {
			err = db.Close()
			if err != nil {
				logger.Warnf("Could not close database: %s", err)
			}
		}

		select {
		case <-stopChan:
			logger.Info("Got stop signal. Finishing in-flight work.")
			cancel()
			cleanup()
			logger.Info("Work done. Goodbye!")
		case <-doneChan:
			cleanup()
			logger.Info("Done scraping. Goodbye!")
		}
	},
}

func init() {
	// scraper
	runCmd.Flags().StringArray("scraper.tokens", []string{}, "The list of tokens the scraper can use")
	_ = viper.BindPFlag("scraper.tokens", runCmd.Flag("scraper.tokens"))

	runCmd.Flags().StringArray("scraper.folders", []string{}, "The list of folders to start from")
	_ = viper.BindPFlag("scraper.folders", runCmd.Flag("scraper.folders"))

	runCmd.Flags().Int("scraper.rps", 0, "Number of request / second / token")
	_ = viper.BindPFlag("scraper.rps", runCmd.Flag("scraper.rps"))
}
