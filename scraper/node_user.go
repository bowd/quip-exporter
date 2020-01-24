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

type UserNode struct {
	*BaseNode
	user *types.QuipUser
}

func NewUserNode(ctx context.Context, id string) interfaces.INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &UserNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", types.NodeTypes.User).
				WithField("id", id),
			id:  id,
			wg:  wg,
			ctx: ctx,
		},
	}
}

func (node *UserNode) Type() types.NodeType {
	return types.NodeTypes.User
}

func (node *UserNode) ID() string {
	return node.id
}

func (node *UserNode) Path() string {
	return path.Join("data", "users", node.id+".json")
}

func (node *UserNode) Children() []interfaces.INode {
	return []interfaces.INode{}
}

func (node *UserNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient) error {
	if node.ctx.Err() != nil {
		return nil
	}

	var user *types.QuipUser
	var err error
	user, err = repo.GetUser(node)
	if err != nil && repositories.IsNotFoundError(err) {
		user, err = quip.GetUser(node.id)
		if err != nil && client.IsUnauthorizedError(err) {
			node.logger.Warn("skipping unauthorized")
			return nil
		} else if err != nil {
			node.logger.Errorln(err)
			return err
		}
		if err := repo.SaveNodeJSON(node, user); err != nil {
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
