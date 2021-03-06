package client

import (
	"encoding/json"
	"fmt"
	"github.com/bowd/quip-exporter/types"
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
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	return qc.getMap(qc.url.Folders(ids), token)
}

func (qc *QuipClient) getThreads(ids []string) (map[string][]byte, error) {
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	return qc.getMap(qc.url.Threads(ids), token)
}

func (qc *QuipClient) getUsers(ids []string) (map[string][]byte, error) {
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	return qc.getMap(qc.url.Users(ids), token)
}

func (qc *QuipClient) getCurrentUser() ([]byte, error) {
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	return qc.getBytes(qc.url.CurrentUser(), token)
}

func (qc *QuipClient) getThreadComments(threadID string) ([]*types.QuipMessage, error) {
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	allComments := make([]*types.QuipMessage, 0, 10)
	var cursor *uint64
	for {
		data, err := qc.getBytes(qc.url.ThreadComments(threadID, cursor), token)
		if err != nil {
			return nil, err
		}
		var comments []*types.QuipMessage
		err = json.Unmarshal(data, &comments)
		if err != nil {
			return nil, err
		}
		if len(comments) == 0 {
			break
		}
		allComments = append(allComments, comments...)
		nextCursor := comments[len(comments)-1].CreatedUsec - 1
		cursor = &nextCursor
		<-time.After(time.Duration(int(1000/qc.rps) * int(time.Millisecond)))
	}
	return allComments, nil
}

func (qc *QuipClient) exportThread(threadID string, exportType string) ([]byte, error) {
	token := qc.checkoutToken()
	defer qc.checkinToken(token)
	return qc.getBytes(qc.url.ExportThread(threadID, exportType), token)
}

func (qc *QuipClient) getBlob(threadID, blobID string) ([]byte, error) {
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	return qc.getBytes(qc.url.Blob(threadID, blobID), token)
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
	if resp.StatusCode == 404 {
		return nil, DeletedError{}
	}
	return []byte(rawBody), nil
}
