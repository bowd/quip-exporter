package scraper

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"regexp"
)

type ThreadHTMLNode struct {
	*ThreadNode
}

var blobRegexp *regexp.Regexp = regexp.MustCompile("blob/([0-9a-zA-Z]+)/([0-9a-zA-z]+)")

func NewThreadHTMLNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadHTMLNode{
		ThreadNode: &ThreadNode{
			Node: &Node{
				path: parent.path,
				id:   parent.id,
				ctx:  ctx,
				wg:   wg,
			},
			thread: parent.thread,
		},
	}
}

func (node *ThreadHTMLNode) ID() string {
	return fmt.Sprintf("thread:%s:html [%s]", node.id, node.path)
}

func (node *ThreadHTMLNode) Children() []INode {
	matches := blobRegexp.FindAllStringSubmatch(node.thread.HTML, -1)
	children := make([]INode, 0, 0)
	for _, match := range matches {
		if match[1] == node.id {
			children = append(children, NewBlobNode(node.ThreadNode, match[2]))
		}
	}
	return []INode{}
}

func (node *ThreadHTMLNode) Process(scraper *Scraper) error {
	scraper.logger.Debugf("Handling thread:%s:HTML [%s/%s]", node.id, node.path, node.thread.Filename())
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := scraper.repo.HasExportedHTML(node.id)
	if err != nil {
		return err
	}

	if !isExported {
		return scraper.repo.SaveThreadHTML(node.path, node.thread)
	}
	return nil
}
