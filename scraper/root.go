package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
	"golang.org/x/sync/errgroup"
)

type RootNode struct {
	*Node
	currentUser *types.QuipUser
}

func NewRootNode(ctx context.Context) INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &RootNode{
		Node: &Node{
			path: "/",
			wg:   wg,
			ctx:  ctx,
		},
	}
}

func (node *RootNode) ID() string {
	return "root"
}

func (node *RootNode) Children() []INode {
	children := make([]INode, 0, 0)
	for _, folderID := range node.currentUser.Folders() {
		children = append(children, NewFolderNode(node.ctx, node.path, folderID))
	}

	return children
}

func (node *RootNode) Process(scraper *Scraper) error {
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
		scraper.logger.Debugf("Loaded current user from repository")
	}
	node.currentUser = currentUser
	return nil
}
