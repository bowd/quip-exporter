package repositories

import (
	"encoding/json"
	"github.com/allegro/bigcache"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"reflect"
)

type CacheRepository struct {
	base   interfaces.IRepository
	cache  *bigcache.BigCache
	logger *logrus.Entry
}

func NewCacheRepository(repo interfaces.IRepository, config bigcache.Config) (interfaces.IRepository, error) {
	cache, err := bigcache.NewBigCache(config)
	return &CacheRepository{
		base:   repo,
		cache:  cache,
		logger: logrus.WithField("module", "cache-repository"),
	}, err
}

type getter = func(interfaces.INode) (interface{}, error)

func (repo *CacheRepository) get(node interfaces.INode, model interface{}, get getter) error {
	data, err := repo.cache.Get(node.Path())
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			fromBase, err := get(node)
			if err != nil {
				return err
			}
			if err := repo.save(node, fromBase); err != nil {
				repo.logger.Warn("Could not write to cache: ", err)
			}
			reflect.ValueOf(model).Elem().Set(reflect.Indirect(reflect.ValueOf(fromBase)))
			return nil
		}
		return err
	} else {
		return json.Unmarshal(data, model)
	}
}

func (repo *CacheRepository) save(node interfaces.INode, model interface{}) error {
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return repo.cache.Set(node.Path(), data)
}

func (repo *CacheRepository) getCurrentUserFromBase(node interfaces.INode) (interface{}, error) {
	return repo.base.GetCurrentUser(node)
}
func (repo *CacheRepository) GetCurrentUser(node interfaces.INode) (*types.QuipUser, error) {
	var model types.QuipUser
	if err := repo.get(node, &model, repo.getCurrentUserFromBase); err != nil {
		return nil, err
	}
	return &model, nil
}

func (repo *CacheRepository) getUserFromBase(node interfaces.INode) (interface{}, error) {
	return repo.base.GetUser(node)
}
func (repo *CacheRepository) GetUser(node interfaces.INode) (*types.QuipUser, error) {
	var model types.QuipUser
	if err := repo.get(node, &model, repo.getUserFromBase); err != nil {
		return nil, err
	}
	return &model, nil
}

func (repo *CacheRepository) getFolderFromBase(node interfaces.INode) (interface{}, error) {
	return repo.base.GetFolder(node)
}
func (repo *CacheRepository) GetFolder(node interfaces.INode) (*types.QuipFolder, error) {
	var model types.QuipFolder
	if err := repo.get(node, &model, repo.getFolderFromBase); err != nil {
		return nil, err
	}
	return &model, nil
}

func (repo *CacheRepository) getThreadFromBase(node interfaces.INode) (interface{}, error) {
	return repo.base.GetThread(node)
}
func (repo *CacheRepository) GetThread(node interfaces.INode) (*types.QuipThread, error) {
	var model types.QuipThread
	if err := repo.get(node, &model, repo.getThreadFromBase); err != nil {
		return nil, err
	}
	return &model, nil
}

func (repo *CacheRepository) getThreadCommentsFromBase(node interfaces.INode) (interface{}, error) {
	return repo.base.GetThreadComments(node)
}
func (repo *CacheRepository) GetThreadComments(node interfaces.INode) ([]*types.QuipMessage, error) {
	var model []*types.QuipMessage
	if err := repo.get(node, &model, repo.getThreadCommentsFromBase); err != nil {
		return nil, err
	}
	return model, nil
}

func (repo *CacheRepository) NodeExists(node interfaces.INode) (bool, error) {
	if _, err := repo.cache.Get(node.Path()); err != nil {
		return repo.base.NodeExists(node)
	} else {
		return true, nil
	}
}

func (repo *CacheRepository) SaveNodeJSON(node interfaces.INode, data interface{}) error {
	return repo.base.SaveNodeJSON(node, data)
}

func (repo *CacheRepository) SaveNodeRaw(node interfaces.INode, data []byte) error {
	return repo.base.SaveNodeRaw(node, data)
}

func (repo *CacheRepository) MakeArchiveCopy(source interfaces.INode, dest interfaces.INode) error {
	return repo.base.MakeArchiveCopy(source, dest)
}
