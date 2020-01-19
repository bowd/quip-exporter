package client

import (
	"encoding/json"
	"fmt"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type QuipClient struct {
	tokens          []string
	logger          *logrus.Entry
	rps             int
	batchWait       time.Duration
	maxItemsInBatch int

	folder batchWithLock
	thread batchWithLock

	tokenIndexMutex sync.Mutex
	lastTokenIndex  int
	tokenMutex      map[string]*sync.Mutex
}

type batchWithLock struct {
	mutex *sync.Mutex
	batch *batch
}

type batch struct {
	ids       []string
	data      map[string][]byte
	error     error
	closing   bool
	done      chan struct{}
	batchType BatchType
}

type BatchType = string

const (
	FolderBatch BatchType = "FolderBatch"
	ThreadBatch BatchType = "ThreadBarch"
)

func New(tokens []string, rps int, batchWait time.Duration, maxItemsInBatch int) (*QuipClient, error) {
	qc := &QuipClient{
		logger:          logrus.WithField("module", "quip-client"),
		rps:             rps,
		batchWait:       batchWait,
		maxItemsInBatch: maxItemsInBatch,
		lastTokenIndex:  0,
		tokenMutex:      make(map[string]*sync.Mutex),
		folder: batchWithLock{
			mutex: &sync.Mutex{},
			batch: nil,
		},
		thread: batchWithLock{
			mutex: &sync.Mutex{},
			batch: nil,
		},
	}

	qc.tokens = make([]string, 0, len(tokens))
	for _, token := range tokens {
		if err := qc.testToken(token); err == nil {
			qc.tokens = append(qc.tokens, token)
			qc.tokenMutex[token] = &sync.Mutex{}
		} else {
			qc.logger.Warnf("Skipping token: %s: %s", token, err)
		}
	}

	if len(qc.tokens) == 0 {
		return nil, fmt.Errorf("could not find any valid token in config")
	}

	qc.logger.Debugf("Setup client with %d valid tokens", len(qc.tokens))
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
