package client

import (
	"encoding/json"
	"fmt"
	"github.com/bowd/quip-exporter/types"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type QuipClient struct {
	token           string
	companyID       string
	logger          *logrus.Entry
	rps             int
	batchWait       time.Duration
	maxItemsInBatch int
	url             URLProvider

	folder batchWithLock
	thread batchWithLock
	user   batchWithLock

	tokenConcurrency int
	lastTokenIndex   int
	tokenIndexMutex  sync.Mutex
	tokenMutex       []*sync.Mutex
}

type batchWithLock struct {
	mutex *sync.Mutex
	batch *batch
}

type batch struct {
	ids       []string
	data      map[string][]byte
	error     map[string]error
	closing   bool
	done      chan struct{}
	batchType BatchType
}

type BatchType = string

const (
	FolderBatch BatchType = "FolderBatch"
	ThreadBatch BatchType = "ThreadBatch"
	UserBatch   BatchType = "UserBatch"
)

func New(token, urlProvider, companyID, baseURL string, tokenConcurrency, rps int, batchWait time.Duration, maxItemsInBatch int) (*QuipClient, error) {
	qc := &QuipClient{
		token:            token,
		companyID:        companyID,
		tokenConcurrency: tokenConcurrency,
		logger:           logrus.WithField("module", "quip-client"),
		rps:              rps,
		batchWait:        batchWait,
		maxItemsInBatch:  maxItemsInBatch,
		lastTokenIndex:   0,
		tokenMutex:       make([]*sync.Mutex, tokenConcurrency),
		folder: batchWithLock{
			mutex: &sync.Mutex{},
			batch: nil,
		},
		thread: batchWithLock{
			mutex: &sync.Mutex{},
			batch: nil,
		},
		user: batchWithLock{
			mutex: &sync.Mutex{},
			batch: nil,
		},
	}

	if urlProvider == "default" {
		qc.url = NewDefaultURLProvider(baseURL)
	} else if urlProvider == "admin" {
		qc.url = NewAdminURLProvider(baseURL, companyID)
	}

	if err := qc.testToken(); err != nil {
		return nil, fmt.Errorf("provided token is invalid")
	}

	for i := 0; i < tokenConcurrency; i++ {
		qc.tokenMutex[i] = &sync.Mutex{}
	}

	return qc, nil
}

func (qc *QuipClient) GetFolder(folderID string) (*types.QuipFolder, error) {
	data, err := qc.getFolderThunk(folderID)()
	if err != nil {
		return nil, err
	}

	var folder types.QuipFolder
	err = json.Unmarshal(data, &folder)
	if err != nil {
		qc.logger.Errorln(err)
		qc.logger.Debugln(string(data))
		return nil, err
	}

	return &folder, nil
}

func (qc *QuipClient) GetThread(threadID string) (*types.QuipThread, error) {
	data, err := qc.getThreadThunk(threadID)()
	if err != nil {
		return nil, err
	}

	var thread types.QuipThread
	err = json.Unmarshal(data, &thread)
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (qc *QuipClient) GetCurrentUser() (*types.QuipUser, error) {
	data, err := qc.getCurrentUser()
	if err != nil {
		return nil, err
	}
	var user types.QuipUser
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (qc *QuipClient) GetUser(userID string) (*types.QuipUser, error) {
	data, err := qc.getUserThunk(userID)()
	if err != nil {
		return nil, err
	}
	var user types.QuipUser
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (qc *QuipClient) ExportThreadSlides(threadID string) ([]byte, error) {
	return qc.exportThread(threadID, "pdf")
}

func (qc *QuipClient) ExportThreadDocument(threadID string) ([]byte, error) {
	return qc.exportThread(threadID, "docx")
}

func (qc *QuipClient) ExportThreadSpreadsheet(threadID string) ([]byte, error) {
	return qc.exportThread(threadID, "xlsx")
}

func (qc *QuipClient) GetThreadComments(threadID string) ([]*types.QuipMessage, error) {
	return qc.getThreadComments(threadID)
}

func (qc *QuipClient) GetBlob(threadID, blobID string) ([]byte, error) {
	return qc.getBlob(threadID, blobID)
}

func (qc *QuipClient) ExportUserPhoto(url string) ([]byte, error) {
	resp, data, errs := gorequest.New().Get(url).Retry(
		5,
		250*time.Millisecond,
		http.StatusBadRequest,
		http.StatusInternalServerError,
	).End()
	if errs != nil {
		qc.logger.Debug(resp)
		qc.logger.Error(data)
		qc.logger.Error(errs)
		return nil, errs[0]
	}
	if resp == nil {
		qc.logger.Debug(resp)
		qc.logger.Error(data)
		qc.logger.Error(errs)
		return nil, fmt.Errorf("response is nil")
	}

	return []byte(data), nil
}
