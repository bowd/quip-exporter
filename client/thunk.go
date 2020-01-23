package client

func (qc *QuipClient) getFolderThunk(folderID string) func() ([]byte, error) {
	b := &qc.folder
	return qc.getThunk(b, FolderBatch, folderID)
}

func (qc *QuipClient) getThreadThunk(threadID string) func() ([]byte, error) {
	b := &qc.thread
	return qc.getThunk(b, ThreadBatch, threadID)
}

func (qc *QuipClient) getUserThunk(userID string) func() ([]byte, error) {
	b := &qc.user
	return qc.getThunk(b, UserBatch, userID)
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
	_ = qc.addToBatch(b, id)
	b.mutex.Unlock()

	return func() ([]byte, error) {
		<-bt.done

		data := bt.data[id]
		err := bt.error[id]

		if (data == nil || len(data) == 0) && err == nil {
			qc.logger.Warnln("Response is empty and there's no error:")
			qc.logger.Warnln(bt.batchType, bt.ids)
		}

		return data, err
	}
}
