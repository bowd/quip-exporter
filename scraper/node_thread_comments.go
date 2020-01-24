package scraper

import (
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type ThreadCommentsNode struct {
	*ThreadNode
	comments []*types.QuipMessage
}

func NewThreadCommentsNode(parent *ThreadNode) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadCommentsNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", types.NodeTypes.ThreadComments).
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

func (node ThreadCommentsNode) Type() types.NodeType {
	return types.NodeTypes.ThreadComments
}

func (node ThreadCommentsNode) ID() string {
	return node.id + "/comments"
}

func (node ThreadCommentsNode) Path() string {
	return path.Join("data", "comments", node.id+".json")
}

func (node *ThreadCommentsNode) Children() []interfaces.INode {
	children := make([]interfaces.INode, 0, 10)
	for _, message := range node.comments {
		children = append(
			children,
			NewUserNode(node.ctx, message.AuthorID),
		)
	}
	return children
}

func (node *ThreadCommentsNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient) error {
	if node.ctx.Err() != nil {
		return nil
	}

	comments, err := repo.GetThreadComments(node)
	if err != nil && repositories.IsNotFoundError(err) {
		comments, err = quip.GetThreadComments(node.id)
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
		if err := repo.SaveNodeJSON(node, comments); err != nil {
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
