package types

import "github.com/kennygrant/sanitize"

type QuipFolder struct {
	Folder    *FolderInfo    `json:"folder"`
	MemberIDs []string       `json:"member_ids"`
	Children  []*FolderChild `json:"children"`
}

type FolderInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	CreatedUsec uint64 `json:"created_usec"`
	CreatorID   string `json:"creator_id"`
	Color       string `json:"color"`
	ParentID    string `json:"parent_id"`
}

type FolderChild struct {
	ThreadID *string `json:"thread_id"`
	FolderID *string `json:"folder_id"`
}

func (fc FolderChild) IsThread() bool {
	return fc.ThreadID != nil
}

func (fc FolderChild) IsFolder() bool {
	return fc.FolderID != nil
}

func (qf QuipFolder) PathSegment() string {
	return sanitize.Path(qf.Folder.Title)
}
