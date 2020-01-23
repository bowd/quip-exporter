package client

import (
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"net/http"
	"time"
)

type batchFetcher = func([]string) (map[string][]byte, error)

func (qc *QuipClient) fetchBatch(b *batch) {
	var fetcher batchFetcher
	if b.batchType == FolderBatch {
		fetcher = qc.getFolders
	} else if b.batchType == ThreadBatch {
		fetcher = qc.getThreads
	} else if b.batchType == UserBatch {
		fetcher = qc.getUsers
	}

	b.data, b.error = qc.fetchIds(b.ids, fetcher)
	close(b.done)
}

func (qc *QuipClient) fetchIds(ids []string, fetcher batchFetcher) (map[string][]byte, map[string]error) {
	if len(ids) == 0 {
		return map[string][]byte{}, map[string]error{}
	}

	data, err := fetcher(ids)
	if err == nil {
		return data, map[string]error{}
	}

	if len(ids) == 1 {
		return data, map[string]error{ids[0]: err}
	}

	if IsUnauthorizedError(err) {
		// If it's unauthorised divide and conquer ids
		midPoint := len(ids) / 2
		data1, errors1 := qc.fetchIds(ids[0:midPoint], fetcher)
		data2, errors2 := qc.fetchIds(ids[midPoint:], fetcher)
		err := mergo.Merge(&data1, data2)
		if err != nil {
			panic(err)
		}
		err = mergo.Merge(&errors1, errors2)
		if err != nil {
			panic(err)
		}
		return data1, errors1
	}

	errors := make(map[string]error)
	for _, id := range ids {
		errors[id] = err
	}
	return map[string][]byte{}, errors
}

func (qc *QuipClient) getFolders(ids []string) (map[string][]byte, error) {
	qc.logger.Debug("Waiting for token")
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	qc.logger.Infof("Querying folder batch [%d folders]", len(ids))
	return qc.getMap(batchURL(FOLDERS_MASK, ids), token)
}

func (qc *QuipClient) getThreads(ids []string) (map[string][]byte, error) {
	qc.logger.Debug("Waiting for token")
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	qc.logger.Debugf("Querying threads batch [%d threads]", len(ids))
	return qc.getMap(batchURL(THREADS_MASK, ids), token)
}

func (qc *QuipClient) getUsers(ids []string) (map[string][]byte, error) {
	qc.logger.Debug("Waiting for token")
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	qc.logger.Debugf("Querying users batch [%d users]", len(ids))
	return qc.getMap(batchURL(USERS_MASK, ids), token)
}

func (qc *QuipClient) getCurrentUser() ([]byte, error) {
	qc.logger.Debug("Waiting for token")
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	qc.logger.Debugf("Querying current user")
	return qc.getBytes(currentUserURL(), token)
}

func (qc *QuipClient) getThreadComments(threadID string, cursor *uint64) ([]byte, error) {
	qc.logger.Debug("Waiting for token")
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	qc.logger.Debugf("Querying comments for thread:%s", threadID)
	return qc.getBytes(threadCommentsURL(threadID, cursor), token)
}

func (qc *QuipClient) exportThread(threadID string, exportType string) ([]byte, error) {
	qc.logger.Debug("Waiting for token")
	token := qc.checkoutToken()
	defer qc.checkinToken(token)
	qc.logger.Debugf("Exporting thread %s as %s", threadID, exportType)
	return qc.getBytes(exportThreadURL(threadID, exportType), token)
}

func (qc *QuipClient) getBlob(threadID, blobID string) ([]byte, error) {
	qc.logger.Debug("Waiting for token")
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	qc.logger.Debugf("Querying blob %s:%s", threadID, blobID)
	return qc.getBytes(blobURL(threadID, blobID), token)
}

func (qc *QuipClient) getMap(url string, token Token) (map[string][]byte, error) {
	rawBody, err := qc.getBytes(url, token)
	if err != nil {
		return nil, err
	}
	var body map[string]interface{}
	err = json.Unmarshal([]byte(rawBody), &body)
	if err != nil {
		qc.logger.Error(rawBody)
		return nil, err
	}
	result := make(map[string][]byte)
	for key, value := range body {
		result[key], _ = json.Marshal(value)
	}
	return result, nil
}

func (qc *QuipClient) getBytes(url string, token Token) ([]byte, error) {
	resp, rawBody, errs := qc.getWithToken(url, token).Retry(
		5,
		250*time.Millisecond,
		http.StatusBadRequest,
		http.StatusInternalServerError,
	).End()

	if errs != nil {
		qc.logger.Debug(resp)
		qc.logger.Error(rawBody)
		qc.logger.Error(errs)
		return nil, errs[0]
	}
	if resp == nil {
		qc.logger.Debug(resp)
		qc.logger.Error(rawBody)
		qc.logger.Error(errs)
		return nil, fmt.Errorf("response is nil")
	}
	if resp.StatusCode == 403 {
		return nil, UnauthorizedError{}
	}
	if resp.StatusCode == 503 {
		return nil, RateLimitError{}
	}
	return []byte(rawBody), nil
}
