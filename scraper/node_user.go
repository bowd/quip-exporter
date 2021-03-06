package scraper

import (
	"context"
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
	if node.user != nil && node.user.ProfilePictureURL != nil {
		return []interfaces.INode{
			NewUserPictureNode(node),
		}
	}
	return []interfaces.INode{}
}

func (node *UserNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	if node.ctx.Err() != nil {
		return nil
	}

	var user *types.QuipUser
	var err error
	user, err = repo.GetUser(node)
	if err != nil && repositories.IsNotFoundError(err) {
		if user, err = quip.GetUser(node.id); err != nil {
			return err
		}
		if err := repo.SaveNodeJSON(node, user); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		node.logger.Debug("loaded from repo")
	}
	node.user = user
	return nil
}
