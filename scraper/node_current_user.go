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
	onlyPrivate bool
}

func NewCurrentUserNode(ctx context.Context, onlyPrivate bool) interfaces.INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &CurrentUserNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", types.NodeTypes.CurrentUser),
			path:   "/",
			wg:     wg,
			ctx:    ctx,
		},
		onlyPrivate: onlyPrivate,
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
	for _, folderID := range node.currentUser.Folders(node.onlyPrivate) {
		children = append(children, NewFolderNode(node.ctx, node.path, folderID))
	}

	return children
}

func (node *CurrentUserNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	var currentUser *types.QuipUser
	var err error
	currentUser, err = repo.GetCurrentUser(node)
	if err != nil && repositories.IsNotFoundError(err) {
		if currentUser, err = quip.GetCurrentUser(); err != nil {
			return err
		}
		if err := repo.SaveNodeJSON(node, currentUser); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		node.logger.Debugf("loaded from repository")
	}
	node.currentUser = currentUser
	return nil
}
