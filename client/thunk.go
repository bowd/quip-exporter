package client

func (qc *QuipClient) getFolderThunk(folderID string) func() ([]byte, error) {
	b := &qc.folder
	return qc.getThunk(b, FolderBatch, folderID)
}

func (qc *QuipClient) getThreadThunk(threadID string) func() ([]byte, error) {
	b := &qc.thread
	return qc.getThunk(b, ThreadBatch, threadID)
}

func (qc *QuipClient) getThunk(b *batchWithLock, batchType BatchType, id string) func() ([]byte, error) {
	b.mutex.Lock()
	if b.batch == nil {
		b.batch = &batch{
			done:      make(chan struct{}),
			batchType: batchType,
		}
	}

	bt := b.batch
	pos := qc.addToBatch(b, id)
	b.mutex.Unlock()

	return func() ([]byte, error) {
		<-bt.done

		var data []byte
		if pos < len(bt.data) {
			data = bt.data[id]
		}

		return data, bt.error
	}
}
