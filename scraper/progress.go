package scraper

import (
	"fmt"
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
			NodeTypes.User,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[NodeTypes.User],
				scraper.progress.queued[NodeTypes.User],
			),
		).
		WithField(
			NodeTypes.Folder,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[NodeTypes.Folder],
				scraper.progress.queued[NodeTypes.Folder],
			),
		).
		WithField(
			NodeTypes.Thread,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[NodeTypes.Thread],
				scraper.progress.queued[NodeTypes.Thread],
			),
		).
		WithField(
			NodeTypes.Blob,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[NodeTypes.Blob],
				scraper.progress.queued[NodeTypes.Blob],
			),
		).
		WithField(
			NodeTypes.Blob,
			fmt.Sprintf(
				"%d/%d",
				scraper.progress.done[NodeTypes.Blob],
				scraper.progress.queued[NodeTypes.Blob],
			),
		).
		Infof("progress")
	go scraper.printProgress()
}

func (scraper *Scraper) incrementQueued(node INode) {
	scraper.progressMutex.Lock()
	defer scraper.progressMutex.Unlock()
	total, _ := scraper.progress.queued["total"]
	scraper.progress.queued["total"] = total + 1
	ofType, _ := scraper.progress.queued[node.Type()]
	scraper.progress.queued[node.Type()] = ofType + 1
}

func (scraper *Scraper) incrementDone(node INode) {
	scraper.progressMutex.Lock()
	defer scraper.progressMutex.Unlock()
	total, _ := scraper.progress.done["total"]
	scraper.progress.done["total"] = total + 1
	ofType, _ := scraper.progress.done[node.Type()]
	scraper.progress.done[node.Type()] = ofType + 1
}
