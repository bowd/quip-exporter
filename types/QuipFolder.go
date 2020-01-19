package types

type QuipFolder struct {
	Folder    *FolderInfo
	MemberIDs []string `json:"member_ids"`
	Children  []*FolderChild
}

type FolderInfo struct {
	ID          string
	Title       string
	CreatedUsec uint64 `json:"created_usec"`
	CreatorID   string `json:"creator_id"`
	Color       string
	ParentID    string
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
