package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type CurrentUserNode struct {
	*BaseNode
	currentUser *types.QuipUser
}

func NewCurrentUserNode(ctx context.Context) INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &CurrentUserNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", NodeTypes.CurrentUser),
			path:   "/",
			wg:     wg,
			ctx:    ctx,
		},
	}
}

func (node CurrentUserNode) Type() NodeType {
	return NodeTypes.CurrentUser
}

func (node CurrentUserNode) Children() []INode {
	children := make([]INode, 0, 0)
	for _, folderID := range node.currentUser.Folders() {
		children = append(children, NewFolderNode(node.ctx, node.path, folderID))
	}

	return children
}

func (node *CurrentUserNode) Process(scraper *Scraper) error {
	var currentUser *types.QuipUser
	var err error
	currentUser, err = scraper.repo.GetCurrentUser()
	if err != nil && repositories.IsNotFoundError(err) {
		currentUser, err = scraper.client.GetCurrentUser()
		if err != nil {
			return err
		}
		if err := scraper.repo.SaveCurrentUser(currentUser); err != nil {
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
