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
	progress      struct {
		queued map[string]int
		done   map[string]int
	}
}

func New(client interfaces.IQuipClient, repo interfaces.IRepository) *Scraper {
	return &Scraper{
		logger:        logrus.WithField("module", "quip-scraper"),
		client:        client,
		repo:          repo,
		progressMutex: &sync.Mutex{},
		progress: struct {
			queued map[string]int
			done   map[string]int
		}{
			queued: make(map[string]int),
			done:   make(map[string]int),
		},
	}
}

func (scraper *Scraper) Run(ctx context.Context, done chan bool) {
	go scraper.printProgress()
	root := NewCurrentUserNode(ctx)
	err := scraper.scrape(ctx, root)
	if err != nil {
		scraper.logger.Errorln(err)
	}
	done <- true

}

func (scraper *Scraper) scrape(ctx context.Context, node INode) error {
	err := node.Process(scraper)
	if err != nil {
		return err
	}
	for _, child := range node.Children() {
		scraper.queue(ctx, node, child)
	}

	err = node.Wait()
	go scraper.incrementDone(node)
	return err
}

func (scraper *Scraper) queue(ctx context.Context, parent, child INode) {
	go scraper.incrementQueued(child)
	parent.Go(func() error { return scraper.scrape(ctx, child) })
}
