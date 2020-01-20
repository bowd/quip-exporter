package client

import (
	"time"
)

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (qc *QuipClient) addToBatch(b *batchWithLock, id string) int {
	for i, queuedID := range b.batch.ids {
		if queuedID == id {
			return i
		}
	}

	pos := len(b.batch.ids)
	b.batch.ids = append(b.batch.ids, id)
	if pos == 0 {
		go qc.startTimer(b, b.batch)
	}

	if qc.maxItemsInBatch != 0 && pos >= qc.maxItemsInBatch-1 {
		if !b.batch.closing {
			bt := b.batch
			bt.closing = true
			b.batch = nil
			go qc.fetchBatch(bt)
		}
	}

	return pos
}

func (qc *QuipClient) startTimer(b *batchWithLock, bt *batch) {
	time.Sleep(qc.batchWait)
	b.mutex.Lock()

	// we must have hit a batch limit and are already finalizing this batch
	if bt.closing {
		b.mutex.Unlock()
		return
	}

	b.batch = nil
	b.mutex.Unlock()
	qc.fetchBatch(bt)
}
