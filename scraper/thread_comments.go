package scraper

import (
	"fmt"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
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
			Node: &Node{
				path: parent.path,
				id:   parent.id,
				ctx:  ctx,
				wg:   wg,
			},
			thread: parent.thread,
		},
	}
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
		if err != nil {
			return err
		}
		if err := scraper.repo.SaveThreadComments(node.id, comments); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		scraper.logger.Debugf("Loaded from repo thread:%s:messages [%s]", node.id, node.path)
	}
	node.comments = comments
	return nil

}
