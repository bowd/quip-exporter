package scraper

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"

	"github.com/bowd/quip-exporter/interfaces"
	"github.com/sirupsen/logrus"
)

type Scraper struct {
	client  interfaces.IQuipClient
	folders []string
	done    chan bool
	wg      *errgroup.Group
	logger  *logrus.Entry
}

type Item struct {
	path string
	id   string
}

func (i Item) child(segment, id string) Item {
	return Item{
		path: i.path + segment + "/",
		id:   id,
	}
}

func New(client interfaces.IQuipClient, folders []string) *Scraper {
	return &Scraper{
		logger:  logrus.WithField("module", "quip-scraper"),
		client:  client,
		folders: folders,
	}

}

func (scraper *Scraper) Run(ctx context.Context, done chan bool) {
	scraper.wg, ctx = errgroup.WithContext(ctx)

	for _, folderID := range scraper.folders {
		scraper.queueFolder(ctx, Item{path: "/", id: folderID})
	}

	err := scraper.wg.Wait()

	if err != nil {
		scraper.logger.Errorln(err)
	}
	done <- true
}

func (scraper *Scraper) queueFolder(ctx context.Context, folder Item) {
	scraper.wg.Go(func() error { return scraper.handleFolder(ctx, folder) })
}

func (scraper *Scraper) queueThread(ctx context.Context, thread Item) {
	scraper.wg.Go(func() error { return scraper.handleThread(ctx, thread) })
}

func (scraper *Scraper) handleFolder(ctx context.Context, item Item) error {
	scraper.logger.Debugf("Handling folder:%s [%s]", item.id, item.path)
	if ctx.Err() != nil {
		return nil
	}

	folder, err := scraper.client.GetFolder(item.id)
	if err != nil {
		return fmt.Errorf("could not get folder: %s", err)
	}

	for _, child := range folder.Children {
		if child.IsFolder() {
			scraper.queueFolder(ctx, item.child(folder.Folder.Title, *child.FolderID))
		} else if child.IsThread() {
			scraper.queueThread(ctx, item.child(folder.Folder.Title, *child.ThreadID))
		}
	}
	return nil
}

func (scraper *Scraper) handleThread(ctx context.Context, item Item) error {
	scraper.logger.Debugf("Handling thread:%s [%s]", item.id, item.path)
	if ctx.Err() != nil {
		return nil

	}

	thread, err := scraper.client.GetThread(item.id)
	if err != nil {
		return err
	}
	dirPath := "./output" + item.path
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			return fmt.Errorf("could not create folder: %s", err)
		}
	}

	filePath := dirPath + thread.Filename() + ".html"
	f, err := os.Create(filePath)
	if err != nil {
		scraper.logger.Errorf("Could not create file for: %s (%s)", thread.Thread.Title, filePath)
		return fmt.Errorf("could not open file: %s: %s", filePath, err)
	}
	_, err = f.WriteString(thread.HTML)
	if err != nil {
		scraper.logger.Errorf("Could not write to file %s (%s)", thread.Thread.Title, filePath)
		return fmt.Errorf("could not write file: %s", err)
	}

	_ = f.Close()

	return nil
}
