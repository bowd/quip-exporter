package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type ThreadNode struct {
	*BaseNode
	thread *types.QuipThread
}

func NewThreadNode(ctx context.Context, path, id string) interfaces.INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &ThreadNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", types.NodeTypes.Thread).
				WithField("id", id).
				WithField("path", path),
			id:   id,
			path: path,
			wg:   wg,
			ctx:  ctx,
		},
	}
}

func (node ThreadNode) Path() string {
	return path.Join("data", "threads", node.id+".json")
}

func (node ThreadNode) Type() types.NodeType {
	return types.NodeTypes.Thread
}

func (node *ThreadNode) Children() []interfaces.INode {
	if node.thread == nil {
		return []interfaces.INode{}
	}
	children := make([]interfaces.INode, 0, 0)
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
		node.logger.Infof("found type channel [!!]")
	}
	return children
}

func (node *ThreadNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient) error {
	if node.ctx.Err() != nil {
		return nil
	}

	var thread *types.QuipThread
	var err error

	thread, err = repo.GetThread(node)
	if err != nil && repositories.IsNotFoundError(err) {
		thread, err = quip.GetThread(node.id)
		if err != nil && client.IsUnauthorizedError(err) {
			node.logger.Warn("skipping unauthorised")
			return nil
		} else if err != nil {
			node.logger.Errorln(err)
			return err
		}
		if err := repo.SaveNodeJSON(node, thread); err != nil {
			node.logger.Errorln(err)
			return err
		}
	} else if err != nil {
		node.logger.Errorln(err)
		return err
	} else {
		node.logger.Debugf("loaded from repository")
	}
	node.thread = thread
	return nil
}
