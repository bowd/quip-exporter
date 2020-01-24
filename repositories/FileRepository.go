package repositories

import (
	"encoding/json"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/bowd/quip-exporter/utils"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type FileRepository struct {
	basePath string
	logger   *logrus.Entry
}

func NewFileRepository(basePath string) interfaces.IRepository {
	return &FileRepository{
		basePath: basePath,
		logger:   logrus.WithField("module", "file-repository"),
	}
}

func (fr *FileRepository) NodeExists(node interfaces.INode) (bool, error) {
	nodePath := path.Join(fr.basePath, node.Path())
	return utils.FileExists(nodePath)
}

func (fr *FileRepository) getJSONNode(node interfaces.INode, data interface{}) error {
	nodePath := path.Join(fr.basePath, node.Path())
	fr.logger.Debug("Loading json from: ", nodePath)
	bytes, err := ioutil.ReadFile(nodePath)
	if os.IsNotExist(err) {
		return NotFoundError{}
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, data); err != nil {
		return err
	}
	return nil
}

func (fr *FileRepository) SaveNodeJSON(node interfaces.INode, data interface{}) error {
	nodePath := path.Join(fr.basePath, node.Path())
	return utils.SaveJSONToFile(nodePath, data)
}

func (fr *FileRepository) SaveNodeRaw(node interfaces.INode, data []byte) error {
	nodePath := path.Join(fr.basePath, node.Path())
	return utils.SaveBytesToFile(nodePath, data)
}

func (fr *FileRepository) GetThread(node interfaces.INode) (*types.QuipThread, error) {
	var data types.QuipThread
	if err := fr.getJSONNode(node, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (fr *FileRepository) GetFolder(node interfaces.INode) (*types.QuipFolder, error) {
	var data types.QuipFolder
	if err := fr.getJSONNode(node, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (fr *FileRepository) GetUser(node interfaces.INode) (*types.QuipUser, error) {
	var data types.QuipUser
	if err := fr.getJSONNode(node, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (fr *FileRepository) GetCurrentUser(node interfaces.INode) (*types.QuipUser, error) {
	var data types.QuipUser
	if err := fr.getJSONNode(node, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (fr *FileRepository) GetThreadComments(node interfaces.INode) ([]*types.QuipMessage, error) {
	var data []*types.QuipMessage
	if err := fr.getJSONNode(node, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (fr *FileRepository) MakeArchiveCopy(source, dest interfaces.INode) error {
	sourcePath := path.Join(fr.basePath, source.Path())
	destPath := path.Join(fr.basePath, dest.Path())
	fr.logger.Infof("%s -> %s", sourcePath, destPath)
	err := utils.EnsureDir(path.Dir(destPath))
	if err != nil {
		return err
	}

	in, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	err = in.Close()
	if err != nil {
		return err
	}
	err = out.Close()
	if err != nil {
		return err
	}
	return nil
}
