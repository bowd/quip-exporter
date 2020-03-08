package commands

import (
	"context"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/scraper"
	"github.com/bowd/quip-exporter/utils"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bowd/quip-exporter/client"
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Start scraper",
	Long:  "Scrape Quip starting from the provided token's current user",
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
		repo := repositories.NewFileRepository(
			viper.GetString("repo.basePath"),
		)
		quipClient, err := client.New(
			viper.GetString("scraper.token"),
			viper.GetString("scraper.company-id"),
			viper.GetInt("scraper.tokenConcurrency"),
			viper.GetInt("scraper.rps"),
			viper.GetDuration("scraper.batch.wait"),
			viper.GetInt("scraper.batch.maxItems"),
		)
		if err != nil {
			logger.Errorln(err)
			return
		}

		scraper := scraper.New(quipClient, repo, viper.GetStringSlice("scraper.blacklist"))
		go scraper.Run(ctx, doneChan)

		cleanup := func() {
			// Cleanup here
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
	scrapeCmd.Flags().StringArray("scraper.tokens", []string{}, "The list of tokens the scraper can use")
	_ = viper.BindPFlag("scraper.tokens", scrapeCmd.Flag("scraper.tokens"))

	scrapeCmd.Flags().StringArray("scraper.folders", []string{}, "The list of folders to start from")
	_ = viper.BindPFlag("scraper.folders", scrapeCmd.Flag("scraper.folders"))

	scrapeCmd.Flags().Int("scraper.rps", 0, "Number of request / second / token")
	_ = viper.BindPFlag("scraper.rps", scrapeCmd.Flag("scraper.rps"))
}
