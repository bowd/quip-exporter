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

type CurrentUserNode struct {
	*BaseNode
	currentUser *types.QuipUser
}

func NewCurrentUserNode(ctx context.Context) interfaces.INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &CurrentUserNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", types.NodeTypes.CurrentUser),
			path:   "/",
			wg:     wg,
			ctx:    ctx,
		},
	}
}

func (node CurrentUserNode) Type() types.NodeType {
	return types.NodeTypes.CurrentUser
}

func (node CurrentUserNode) ID() string {
	return "root"
}

func (node CurrentUserNode) Path() string {
	return path.Join("data", "root.json")
}

func (node CurrentUserNode) Children() []interfaces.INode {
	children := make([]interfaces.INode, 0, 0)
	for _, folderID := range node.currentUser.Folders() {
		children = append(children, NewFolderNode(node.ctx, node.path, folderID))
	}

	return children
}

func (node *CurrentUserNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient) error {
	var currentUser *types.QuipUser
	var err error
	currentUser, err = repo.GetCurrentUser(node)
	if err != nil && repositories.IsNotFoundError(err) {
		currentUser, err = quip.GetCurrentUser()
		if err != nil {
			node.logger.Errorln(err)
			return err
		}
		if err := repo.SaveNodeJSON(node, currentUser); err != nil {
			node.logger.Errorln(err)
			return err
		}
	} else if err != nil {
		node.logger.Errorln(err)
		return err
	} else {
		node.logger.Debugf("loaded from repository")
	}
	node.currentUser = currentUser
	return nil
}
