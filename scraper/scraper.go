package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Scraper struct {
	client  interfaces.IQuipClient
	folders []string
	done    chan bool
	wg      *errgroup.Group
	logger  *logrus.Entry
	repo    interfaces.IRepository
}

func New(client interfaces.IQuipClient, repo interfaces.IRepository, folders []string) *Scraper {
	return &Scraper{
		logger:  logrus.WithField("module", "quip-scraper"),
		client:  client,
		folders: folders,
		repo:    repo,
	}

}

func (scraper *Scraper) Run(ctx context.Context, done chan bool) {
	root := NewRootNode(ctx)
	err := scraper.scrape(ctx, root)
	if err != nil {
		scraper.logger.Errorln(err)
	}
	done <- true

}

func (scraper *Scraper) scrape(ctx context.Context, node INode) error {
	scraper.logger.Debugf("Scraping %s", node.ID())
	err := node.Process(scraper)
	if err != nil {
		return err
	}
	for _, child := range node.Children() {
		scraper.queue(ctx, node, child)
	}

	return node.Wait()
}

func (scraper *Scraper) queue(ctx context.Context, parent, child INode) {
	parent.Go(func() error { return scraper.scrape(ctx, child) })
}
