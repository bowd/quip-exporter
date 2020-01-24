package scraper

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type BaseNode struct {
	path   string
	id     string
	wg     *errgroup.Group
	ctx    context.Context
	logger *logrus.Entry
}

func (node *BaseNode) Go(fn func() error) {
	node.wg.Go(fn)
}

func (node BaseNode) ID() string {
	return node.id
}

func (node *BaseNode) Wait() error {
	return node.wg.Wait()
}
