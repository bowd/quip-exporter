package search

import (
	"github.com/avast/retry-go"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"sync"
)

var logger = logrus.WithField("module", "search")

type Search struct {
	Index        bleve.Index
	batch        *bleve.Batch
	lock         *sync.Mutex
	maxBatchSize int
	ids          map[string]int
}

func New(path string) *Search {
	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		logger.Info("Index not found, building new index")
		indexMapping, err := buildIndexMapping()
		if err != nil {
			logger.Fatal(err)
		}
		index, err = bleve.New(path, indexMapping)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Info("Opening index")
	}

	return &Search{
		Index:        index,
		lock:         &sync.Mutex{},
		maxBatchSize: 100,
		ids:          make(map[string]int),
	}
}

func buildIndexMapping() (mapping.IndexMapping, error) {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	threadMapping := bleve.NewDocumentMapping()
	threadMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	threadMapping.AddFieldMappingsAt("content", englishTextFieldMapping)

	authorMapping := bleve.NewDocumentMapping()
	authorMapping.AddFieldMappingsAt("name", englishTextFieldMapping)

	folderMapping := bleve.NewDocumentMapping()
	folderMapping.AddFieldMappingsAt("name", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("thread", threadMapping)
	indexMapping.AddDocumentMapping("author", authorMapping)
	indexMapping.AddDocumentMapping("folder", folderMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

func (search *Search) IndexFolder(id, title string) error {
	return search.Index.Index(id, map[string]string{
		"type": "folder",
		"name": title,
	})
}

func (search *Search) IndexThread(id, title, html string) error {
	return search.singleIndexThread(id, title, html)
}

func (search *Search) singleIndexThread(id, title, html string) error {
	return search.Index.Index(id, map[string]string{
		"type":    "thread",
		"name":    title,
		"content": html,
	})
}

func (search *Search) IsIndexed(id string) (bool, error) {
	result, err := search.Index.Search(bleve.NewSearchRequest(bleve.NewDocIDQuery([]string{id})))
	if err != nil {
		return false, err
	}
	return result.Hits.Len() > 0, nil
}

func (search *Search) batchIndexThread(id string, title, html string) error {
	search.lock.Lock()
	defer search.lock.Unlock()
	if search.batch == nil {
		search.batch = search.Index.NewBatch()
	}

	err := search.batch.Index(id, map[string]string{
		"type":    "thread",
		"name":    title,
		"content": html,
	})
	search.ids[id] = len(html)
	if err != nil {
		logger.Debug(id)
		return err
	}

	if search.batch.Size() >= search.maxBatchSize {
		logger.Debug("Indexing thread batch")
		logger.Debugf("%d items in batch :: %d", search.batch.Size(), search.batch.TotalDocsSize())
		err := retry.Do(
			func() error {
				return search.Index.Batch(search.batch)
			},
			retry.Attempts(4),
		)
		search.batch = nil
		spew.Dump(search.ids)
		search.ids = make(map[string]int)
		return err
	}
	return nil
}

func (search *Search) Close() error {
	search.lock.Lock()
	defer search.lock.Unlock()

	if search.batch != nil {
		logger.Debugf("%d items in batch", search.batch.Size())
		return search.Index.Batch(search.batch)
	}
	return nil
}
