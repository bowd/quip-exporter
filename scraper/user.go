package scraper

import (
	"context"
	"fmt"
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
	"golang.org/x/sync/errgroup"
)

type UserNode struct {
	*Node
	user *types.QuipUser
}

func NewUserNode(ctx context.Context, id string) INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &UserNode{
		Node: &Node{
			id:  id,
			wg:  wg,
			ctx: ctx,
		},
	}
}

func (node *UserNode) ID() string {
	return fmt.Sprintf("user:%s", node.id)
}

func (node *UserNode) Children() []INode {
	return []INode{}
}

func (node *UserNode) Process(scraper *Scraper) error {
	scraper.logger.Debugf("Handling user:%s", node.id)
	if node.ctx.Err() != nil {
		return nil
	}

	var user *types.QuipUser
	var err error
	user, err = scraper.repo.GetUser(node.id)
	if err != nil && repositories.IsNotFoundError(err) {
		user, err = scraper.client.GetUser(node.id)
		if err != nil && client.IsUnauthorizedError(err) {
			scraper.logger.Debugf("Ignoring Unauthorized user:%s [%s]", node.id, node.path)
			return nil
		} else if err != nil {
			return err
		}
		if err := scraper.repo.SaveUser(user); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		scraper.logger.Debugf("Loaded from repo user:%s", node.id)
	}
	node.user = user
	return nil
}
