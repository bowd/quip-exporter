package scraper

import "time"

func (scraper *Scraper) printProgress() {
	<-time.After(5 * time.Second)
	scraper.logger.Info("==================================")
	scraper.logger.Infof("Progress: %d/%d nodes", scraper.doneNodes, scraper.queuedNodes)
	scraper.logger.Info("==================================")
	go scraper.printProgress()
}

func (scraper *Scraper) incrementQueued() {
	scraper.progressMutex.Lock()
	defer scraper.progressMutex.Unlock()
	scraper.queuedNodes += 1
}

func (scraper *Scraper) incrementDone() {
	scraper.progressMutex.Lock()
	defer scraper.progressMutex.Unlock()
	scraper.doneNodes += 1
}
