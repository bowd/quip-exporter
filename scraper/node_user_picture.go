package scraper

import (
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type UserPictureNode struct {
	*BaseNode
	user   *types.QuipUser
	exists bool
}

func NewUserPictureNode(parent *UserNode) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &UserPictureNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", types.NodeTypes.UserPicture).
				WithField("id", parent.id),
			id:  parent.id,
			wg:  wg,
			ctx: ctx,
		},
		user:   parent.user,
		exists: false,
	}
}

func (node *UserPictureNode) Type() types.NodeType {
	return types.NodeTypes.UserPicture
}

func (node *UserPictureNode) ID() string {
	return node.id
}

func (node *UserPictureNode) Path() string {
	return path.Join("data", "profile_pictures", node.id+".png")
}

func (node *UserPictureNode) Children() []interfaces.INode {
	return []interfaces.INode{}
}

func (node *UserPictureNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient) error {
	if node.ctx.Err() != nil {
		return nil
	}
	if node.user.ProfilePictureURL == nil {
		node.logger.Warn("skipping: url missing")
		return nil
	}
	isExported, err := repo.NodeExists(node)
	if err != nil {
		node.logger.Errorln(err)
		return err
	}

	if !isExported {
		data, err := quip.ExportUserPhoto(*node.user.ProfilePictureURL)
		if err != nil {
			node.logger.Errorln(err)
			return err
		}
		if err := repo.SaveNodeRaw(node, data); err != nil {
			node.logger.Errorln(err)
			return err
		} else {
			node.exists = true
		}
	} else {
		node.exists = true
		node.logger.Debugf("already exported")
	}
	return nil
}
