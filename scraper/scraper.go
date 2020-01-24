package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
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
	seenMap       *sync.Map
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
		seenMap:       &sync.Map{},
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

func (scraper *Scraper) scrape(ctx context.Context, node interfaces.INode) error {
	err := node.Process(scraper.repo, scraper.client)
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

func (scraper *Scraper) queue(ctx context.Context, parent, child interfaces.INode) {
	if !scraper.shouldSkip(child) {
		go scraper.incrementQueued(child)
		parent.Go(func() error { return scraper.scrape(ctx, child) })
	}
}

func (scraper *Scraper) shouldSkip(child interfaces.INode) bool {
	if child.Type() == types.NodeTypes.User {
		key := child.Type() + "::" + child.ID()
		_, seen := scraper.seenMap.Load(key)
		if seen {
			return true
		} else {
			scraper.seenMap.Store(key, true)
			return false
		}
	} else {
		return false
	}

}
