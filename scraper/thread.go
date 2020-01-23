package scraper

import (
	"context"
	"fmt"
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
	"golang.org/x/sync/errgroup"
)

type ThreadNode struct {
	*Node
	thread *types.QuipThread
}

func NewThreadNode(ctx context.Context, path, id string) INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &ThreadNode{
		Node: &Node{
			id:   id,
			path: path,
			wg:   wg,
			ctx:  ctx,
		},
	}
}

func (node *ThreadNode) ID() string {
	return fmt.Sprintf("thread:%s [%s]", node.id, node.path)
}

func (node *ThreadNode) Children() []INode {
	if node.thread == nil {
		return []INode{}
	}
	children := make([]INode, 0, 0)
	children = append(
		children,
		NewThreadCommentsNode(node),
		NewUserNode(node.ctx, node.thread.Thread.AuthorID),
	)

	if !node.thread.IsChannel() {
		children = append(children, NewThreadHTMLNode(node))
	}

	if node.thread.IsSlides() {
		children = append(children, NewThreadSlidesNode(node))
	}
	if node.thread.IsDocument() {
		children = append(children, NewThreadDocumentNode(node))
	}
	if node.thread.IsSpreadsheet() {
		children = append(children, NewThreadSpreadsheetNode(node))
	}
	if node.thread.IsChannel() {
		fmt.Println("!!!!!! Found channel thread: ", node.thread.Thread.ID)
	}
	return children
}

func (node *ThreadNode) Process(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}

	var thread *types.QuipThread
	var err error

	thread, err = scraper.repo.GetThread(node.id)
	if err != nil && repositories.IsNotFoundError(err) {
		thread, err = scraper.client.GetThread(node.id)
		if err != nil && client.IsUnauthorizedError(err) {
			scraper.logger.Debugf("Ignoring Unauthorized thread:%s [%s]", node.id, node.path)
			return nil
		} else if err != nil {
			return err
		}
		if err := scraper.repo.SaveThread(thread); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		scraper.logger.Debugf("Loaded from repo thread:%s [%s]", node.id, node.path)
	}
	node.thread = thread
	return nil
}
