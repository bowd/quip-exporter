package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/types"
	"github.com/bowd/quip-exporter/utils"
	"golang.org/x/sync/errgroup"
	"path"
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

func (node *RootNode) Children() []INode {
	children := make([]INode, 0, 0)
	for _, folderID := range node.currentUser.Folders() {
		children = append(children, NewFolderNode(node.ctx, node.path, folderID))
	}

	return children
}

func (node *RootNode) Load(scraper *Scraper) error {
	currentUser, err := scraper.client.GetCurrentUser()
	if err != nil {
		return err
	}
	node.currentUser = currentUser
	return node.saveData(scraper)
}

func (node *RootNode) saveData(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}
	if err := utils.SaveJSONToFile(path.Join(DATA_FOLDER, "root.json"), node.currentUser); err != nil {
		return err
	}
	return nil
}
