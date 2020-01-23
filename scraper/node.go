package scraper

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type INode interface {
	Go(func() error)
	Wait() error
	Children() []INode
	Process(*Scraper) error
	Type() NodeType
	ID() string
}

type NodeType = string

var NodeTypes = struct {
	CurrentUser       NodeType
	Blob              NodeType
	User              NodeType
	Folder            NodeType
	Thread            NodeType
	ThreadHTML        NodeType
	ThreadComments    NodeType
	ThreadDocument    NodeType
	ThreadSlides      NodeType
	ThreadSpreadsheet NodeType
}{
	CurrentUser:       "current-user",
	Blob:              "blob",
	User:              "user",
	Folder:            "folder",
	Thread:            "thread",
	ThreadHTML:        "thread-html",
	ThreadComments:    "thread-comments",
	ThreadDocument:    "thread-document",
	ThreadSlides:      "thread-slides",
	ThreadSpreadsheet: "thread-spreadsheet",
}

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
