package scraper

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type INode interface {
	Go(func() error)
	Wait() error
	Children() []INode
	Process(*Scraper) error
	ID() string
}

type Node struct {
	path string
	id   string
	wg   *errgroup.Group
	ctx  context.Context
}

func (node *Node) Go(fn func() error) {
	node.wg.Go(fn)
}

func (node *Node) Wait() error {
	return node.wg.Wait()
}
