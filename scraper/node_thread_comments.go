package scraper

import (
	"fmt"
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type ThreadCommentsNode struct {
	*ThreadNode
	comments []*types.QuipMessage
}

func NewThreadCommentsNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadCommentsNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", NodeTypes.ThreadComments).
					WithField("id", parent.id).
					WithField("path", parent.path),
				path: parent.path,
				id:   parent.id,
				ctx:  ctx,
				wg:   wg,
			},
			thread: parent.thread,
		},
	}
}

func (node ThreadCommentsNode) Type() NodeType {
	return NodeTypes.ThreadComments
}

func (node *ThreadCommentsNode) ID() string {
	return fmt.Sprintf("thread:%s:comments [%s]", node.id, node.path)
}

func (node *ThreadCommentsNode) Children() []INode {
	children := make([]INode, 0, 10)
	for _, message := range node.comments {
		children = append(
			children,
			NewUserNode(node.ctx, message.AuthorID),
		)
	}
	return children
}

func (node *ThreadCommentsNode) Process(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}

	comments, err := scraper.repo.GetThreadComments(node.id)
	if err != nil && repositories.IsNotFoundError(err) {
		comments, err = scraper.client.GetThreadComments(node.id)
		if err != nil && client.IsUnauthorizedError(err) {
			node.logger.Warn("skipping unauthorised")
			return nil
		} else if err != nil && client.IsDeletedError(err) {
			node.logger.Warn("skipping deleted")
			return nil
		} else if err != nil {
			node.logger.Errorln(err)
			return err
		}
		if err := scraper.repo.SaveThreadComments(node.id, comments); err != nil {
			node.logger.Errorln(err)
			return err
		}
	} else if err != nil {
		node.logger.Errorln(err)
		return err
	} else {
		node.logger.Debugf("loaded from repo")
	}
	node.comments = comments
	return nil

}
