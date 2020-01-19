package commands

import (
	"context"
	"github.com/bowd/quip-exporter/scraper"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

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
		stopChan := make(chan os.Signal, 1)
		doneChan := make(chan bool, 1)
		signal.Notify(stopChan, syscall.SIGINT)
		signal.Notify(stopChan, syscall.SIGTERM)
		ctx, cancel := context.WithCancel(context.Background())
		quipClient, err := client.New(
			viper.GetStringSlice("scraper.tokens"),
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
			viper.GetStringSlice("scraper.folders"),
		)

		go scraper.Run(ctx, doneChan)

		select {
		case <-stopChan:
			log.Info("Got stop signal. Finishing work.")
			cancel()
			// close whatever there is to close
			log.Info("Work done. Goodbye!")
		case <-doneChan:
			log.Info("Done scraping. Goodbye!")
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
