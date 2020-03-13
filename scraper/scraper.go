package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
)

type Scraper struct {
	blacklist     map[types.NodeType]bool
	done          chan bool
	wg            *errgroup.Group
	logger        *logrus.Entry
	client        interfaces.IQuipClient
	repo          interfaces.IRepository
	search        interfaces.ISearchIndex
	progressMutex *sync.Mutex
	seenMap       *sync.Map
	onlyPrivate   bool
	progress      struct {
		queued map[string]int
		done   map[string]int
	}
}

func New(client interfaces.IQuipClient, repo interfaces.IRepository, search interfaces.ISearchIndex, _blacklist []string, onlyPrivate bool) *Scraper {
	blacklist := make(map[types.NodeType]bool)
	for _, node := range _blacklist {
		blacklist[node] = true
	}
	return &Scraper{
		logger:        logrus.WithField("module", "quip-scraper"),
		client:        client,
		search:        search,
		repo:          repo,
		progressMutex: &sync.Mutex{},
		seenMap:       &sync.Map{},
		blacklist:     blacklist,
		onlyPrivate:   onlyPrivate,
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
	root := NewCurrentUserNode(ctx, scraper.onlyPrivate)
	err := scraper.scrape(ctx, root)
	if err != nil {
		scraper.logger.Errorln(err)
	}
	done <- true
}

func (scraper *Scraper) scrape(ctx context.Context, node interfaces.INode) error {
	defer scraper.incrementDone(node)
	nodeLogger := scraper.logger.WithField("type", node.Type()).WithField("id", node.ID())
	err := node.Process(scraper.repo, scraper.client, scraper.search)
	if err != nil {
		if client.IsUnauthorizedError(err) {
			nodeLogger.Warn("skipping unauthorized")
			return nil
		} else if client.IsDeletedError(err) {
			nodeLogger.Warn("skipping deleted")
			return nil
		} else if strings.Contains(err.Error(), "key too large") {
			nodeLogger.Warn("skipping index (key too large)")
			return nil
		} else {
			nodeLogger.Error(err)
			spew.Dump(node)
			return err
		}
	}

	for _, child := range node.Children() {
		scraper.queue(ctx, node, child)
	}
	err = node.Wait()
	if err != nil {
		nodeLogger.Error(err)
		spew.Dump(node)
		return err
	}
	return err
}

func (scraper *Scraper) queue(ctx context.Context, parent, child interfaces.INode) {
	if !scraper.shouldSkip(child) {
		go scraper.incrementQueued(child)
		parent.Go(func() error { return scraper.scrape(ctx, child) })
	}
}

func (scraper *Scraper) shouldSkip(child interfaces.INode) bool {
	if scraper.blacklist[child.Type()] == true {
		return true
	} else if child.Type() == types.NodeTypes.User {
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
