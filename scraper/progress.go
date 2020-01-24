package scraper

import (
	"fmt"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"time"
)

func (scraper *Scraper) printProgress() {
	<-time.After(5 * time.Second)
	scraper.logger.
		WithField(
			"total",
			fmt.Sprintf("%d/%d", scraper.progress.done["total"], scraper.progress.queued["total"]),
		).
		WithField(
			types.NodeTypes.User,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[types.NodeTypes.User],
				scraper.progress.queued[types.NodeTypes.User],
			),
		).
		WithField(
			types.NodeTypes.Folder,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[types.NodeTypes.Folder],
				scraper.progress.queued[types.NodeTypes.Folder],
			),
		).
		WithField(
			types.NodeTypes.Thread,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[types.NodeTypes.Thread],
				scraper.progress.queued[types.NodeTypes.Thread],
			),
		).
		WithField(
			types.NodeTypes.Blob,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[types.NodeTypes.Blob],
				scraper.progress.queued[types.NodeTypes.Blob],
			),
		).
		WithField(
			types.NodeTypes.Archive,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[types.NodeTypes.Archive],
				scraper.progress.queued[types.NodeTypes.Archive],
			),
		).
		WithField(
			types.NodeTypes.ThreadComments,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[types.NodeTypes.ThreadComments],
				scraper.progress.queued[types.NodeTypes.ThreadComments],
			),
		).
		Infof("progress")
	go scraper.printProgress()
}

func (scraper *Scraper) incrementQueued(node interfaces.INode) {
	scraper.progressMutex.Lock()
	defer scraper.progressMutex.Unlock()
	total, _ := scraper.progress.queued["total"]
	scraper.progress.queued["total"] = total + 1
	ofType, _ := scraper.progress.queued[node.Type()]
	scraper.progress.queued[node.Type()] = ofType + 1
}

func (scraper *Scraper) incrementDone(node interfaces.INode) {
	scraper.progressMutex.Lock()
	defer scraper.progressMutex.Unlock()
	total, _ := scraper.progress.done["total"]
	scraper.progress.done["total"] = total + 1
	ofType, _ := scraper.progress.done[node.Type()]
	scraper.progress.done[node.Type()] = ofType + 1
}
