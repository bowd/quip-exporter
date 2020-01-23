package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type UserNode struct {
	*BaseNode
	user *types.QuipUser
}

func NewUserNode(ctx context.Context, id string) INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &UserNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", NodeTypes.User).
				WithField("id", id),
			id:  id,
			wg:  wg,
			ctx: ctx,
		},
	}
}

func (node *UserNode) Type() NodeType {
	return NodeTypes.User
}

func (node *UserNode) Children() []INode {
	return []INode{}
}

func (node *UserNode) Process(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}

	var user *types.QuipUser
	var err error
	user, err = scraper.repo.GetUser(node.id)
	if err != nil && repositories.IsNotFoundError(err) {
		user, err = scraper.client.GetUser(node.id)
		if err != nil && client.IsUnauthorizedError(err) {
			node.logger.Warn("skipping unauthorized")
			return nil
		} else if err != nil {
			node.logger.Errorln(err)
			return err
		}
		if err := scraper.repo.SaveUser(user); err != nil {
			node.logger.Errorln(err)
			return err
		}
	} else if err != nil {
		node.logger.Errorln(err)
		return err
	} else {
		node.logger.Debug("loaded from repo")
	}
	node.user = user
	return nil
}
