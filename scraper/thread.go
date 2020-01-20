package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/types"
	"github.com/bowd/quip-exporter/utils"
	"golang.org/x/sync/errgroup"
	"path"
)

type ThreadNode struct {
	*Node
	thread *types.QuipThread
}

func NewThreadNode(ctx context.Context, path, id string) INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &ThreadNode{
		Node: &Node{
			id:   id,
			path: path,
			wg:   wg,
			ctx:  ctx,
		},
	}
}

func (node *ThreadNode) Children() []INode {
	children := make([]INode, 0, 0)
	return children
}

func (node *ThreadNode) Load(scraper *Scraper) error {
	var err error
	scraper.logger.Debugf("Handling thread:%s [%s]", node.id, node.path)
	if node.ctx.Err() != nil {
		return nil
	}

	thread, err := scraper.client.GetThread(node.id)
	node.thread = thread
	if err != nil {
		return err
	}
	err = node.saveHTML(scraper)
	if err != nil {
		return err
	}
	err = node.saveData(scraper)
	if err != nil {
		return err
	}

	return nil
}

func (node *ThreadNode) saveHTML(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}
	err := utils.SaveBytesToFile(
		path.Join(FLAT_HTML_FOLDER, "threads", node.id+".html"),
		[]byte(node.thread.HTML),
	)
	if err != nil {
		return err
	}
	err = utils.SaveBytesToFile(
		path.Join(HTML_FOLDER+node.path, node.thread.Filename()),
		[]byte(node.thread.HTML),
	)
	if err != nil {
		return err
	}
	return nil
}

func (node *ThreadNode) saveData(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}
	if err := utils.SaveJSONToFile(path.Join(DATA_FOLDER, "threads", node.id+".json"), node.thread); err != nil {
		return err
	}
	return nil
}
