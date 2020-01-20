package repositories

import (
	"encoding/json"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/bowd/quip-exporter/utils"
	"io/ioutil"
	"os"
	"path"
)

type FileRepository struct {
	basePath string
}

func NewFileRepository(basePath string) interfaces.IRepository {
	return &FileRepository{basePath}
}

func (fr *FileRepository) GetThread(id string) (*types.QuipThread, error) {
	bytes, err := ioutil.ReadFile(path.Join(fr.basePath, "data", "threads", id+".json"))
	if os.IsNotExist(err) {
		return nil, NotFoundError{}
	}
	if err != nil {
		return nil, err
	}
	var thread types.QuipThread
	if err := json.Unmarshal(bytes, &thread); err != nil {
		return nil, err
	}
	return &thread, nil
}

func (fr *FileRepository) SaveThread(thread *types.QuipThread) error {
	if err := utils.SaveJSONToFile(path.Join(fr.basePath, "data", "threads", thread.Thread.ID+".json"), thread); err != nil {
		return err
	}
	return nil
}

func (fr *FileRepository) GetFolder(id string) (*types.QuipFolder, error) {
	bytes, err := ioutil.ReadFile(path.Join(fr.basePath, "data", "folders", id+".json"))
	if os.IsNotExist(err) {
		return nil, NotFoundError{}
	}
	if err != nil {
		return nil, err
	}
	var folder types.QuipFolder
	if err := json.Unmarshal(bytes, &folder); err != nil {
		return nil, err
	}
	return &folder, nil
}

func (fr *FileRepository) SaveFolder(folder *types.QuipFolder) error {
	if err := utils.SaveJSONToFile(path.Join(fr.basePath, "data", "folders", folder.Folder.ID+".json"), folder); err != nil {
		return err
	}
	return nil
}

func (fr *FileRepository) GetUser(id string) (*types.QuipUser, error) {
	bytes, err := ioutil.ReadFile(path.Join(fr.basePath, "data", "users", id+".json"))
	if os.IsNotExist(err) {
		return nil, NotFoundError{}
	}
	if err != nil {
		return nil, err
	}
	var currentUser types.QuipUser
	if err := json.Unmarshal(bytes, &currentUser); err != nil {
		return nil, err
	}
	return &currentUser, nil
}

func (fr *FileRepository) SaveUser(user *types.QuipUser) error {
	if err := utils.SaveJSONToFile(path.Join(fr.basePath, "data", "users", user.ID+".json"), user); err != nil {
		return err
	}
	return nil
}

func (fr *FileRepository) GetCurrentUser() (*types.QuipUser, error) {
	bytes, err := ioutil.ReadFile(path.Join(fr.basePath, "data", "root.json"))
	if os.IsNotExist(err) {
		return nil, NotFoundError{}
	}
	if err != nil {
		return nil, err
	}
	var currentUser types.QuipUser
	if err := json.Unmarshal(bytes, &currentUser); err != nil {
		return nil, err
	}
	return &currentUser, nil
}

func (fr *FileRepository) SaveCurrentUser(user *types.QuipUser) error {
	if err := utils.SaveJSONToFile(path.Join(fr.basePath, "data", "root.json"), user); err != nil {
		return err
	}
	return nil
}

func (fr *FileRepository) HasExportedHTML(threadID string) (bool, error) {
	htmlPath := path.Join(fr.basePath, "data", "html", threadID+".html")
	return utils.FileExists(htmlPath)
}

func (fr *FileRepository) SaveThreadHTML(nodePath string, thread *types.QuipThread) error {
	filename := path.Join(fr.basePath, "data", "html", thread.Thread.ID+".html")
	err := utils.SaveBytesToFile(filename, []byte(thread.HTML))
	if err != nil {
		return nil
	}
	publicFilename := path.Join(fr.basePath+"/archive"+nodePath, thread.Filename()+".html")
	return utils.SaveBytesToFile(publicFilename, []byte(thread.HTML))
}

func (fr *FileRepository) HasExportedSlides(threadID string) (bool, error) {
	htmlPath := path.Join(fr.basePath, "data", "pdf", threadID+".pdf")
	return utils.FileExists(htmlPath)
}

func (fr *FileRepository) SaveThreadSlides(nodePath string, thread *types.QuipThread, pdf []byte) error {
	pdfPath := path.Join(fr.basePath, "data", "pdf", thread.Thread.ID+".pdf")
	err := utils.SaveBytesToFile(pdfPath, pdf)
	if err != nil {
		return nil
	}
	publicFilename := path.Join(fr.basePath+"/archive"+nodePath, thread.Filename()+".pdf")
	return utils.SaveBytesToFile(publicFilename, pdf)
}

func (fr *FileRepository) HasExportedDocument(threadID string) (bool, error) {
	htmlPath := path.Join(fr.basePath, "data", "docs", threadID+".docx")
	return utils.FileExists(htmlPath)
}

func (fr *FileRepository) SaveThreadDocument(nodePath string, thread *types.QuipThread, doc []byte) error {
	pdfPath := path.Join(fr.basePath, "data", "docs", thread.Thread.ID+".docx")
	err := utils.SaveBytesToFile(pdfPath, doc)
	if err != nil {
		return nil
	}
	publicFilename := path.Join(fr.basePath+"/archive"+nodePath, thread.Filename()+".docx")
	return utils.SaveBytesToFile(publicFilename, doc)
}

func (fr *FileRepository) HasExportedSpreadsheet(threadID string) (bool, error) {
	htmlPath := path.Join(fr.basePath, "data", "xls", threadID+".xlsx")
	return utils.FileExists(htmlPath)
}

func (fr *FileRepository) SaveThreadSpreadsheet(nodePath string, thread *types.QuipThread, xls []byte) error {
	pdfPath := path.Join(fr.basePath, "data", "xls", thread.Thread.ID+".xlsx")
	err := utils.SaveBytesToFile(pdfPath, xls)
	if err != nil {
		return nil
	}
	publicFilename := path.Join(fr.basePath+"/archive"+nodePath, thread.Filename()+".xlsx")
	return utils.SaveBytesToFile(publicFilename, xls)
}

func (fr *FileRepository) GetThreadComments(threadID string) ([]*types.QuipMessage, error) {
	bytes, err := ioutil.ReadFile(path.Join(fr.basePath, "data", "comments", threadID+".json"))
	if os.IsNotExist(err) {
		return nil, NotFoundError{}
	}
	if err != nil {
		return nil, err
	}
	var comments []*types.QuipMessage
	if err := json.Unmarshal(bytes, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (fr *FileRepository) SaveThreadComments(threadID string, comments []*types.QuipMessage) error {
	if err := utils.SaveJSONToFile(path.Join(fr.basePath, "data", "comments", threadID+".json"), comments); err != nil {
		return err
	}
	return nil
}
