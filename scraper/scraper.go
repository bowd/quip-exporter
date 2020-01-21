package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"sync"
)

type Scraper struct {
	client        interfaces.IQuipClient
	done          chan bool
	wg            *errgroup.Group
	logger        *logrus.Entry
	repo          interfaces.IRepository
	progressMutex *sync.Mutex
	queuedNodes   uint
	doneNodes     uint
}

func New(client interfaces.IQuipClient, repo interfaces.IRepository) *Scraper {
	return &Scraper{
		logger:        logrus.WithField("module", "quip-scraper"),
		client:        client,
		repo:          repo,
		progressMutex: &sync.Mutex{},
	}
}

func (scraper *Scraper) Run(ctx context.Context, done chan bool) {
	go scraper.printProgress()
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

	err = node.Wait()
	go scraper.incrementDone()
	return err
}

func (scraper *Scraper) queue(ctx context.Context, parent, child INode) {
	go scraper.incrementQueued()
	parent.Go(func() error { return scraper.scrape(ctx, child) })
}
