package types

type Path = []PathFolder

type PathFolder struct {
	FolderID string `json:"folder_id"`
	Title    string `json:"title"`
}
