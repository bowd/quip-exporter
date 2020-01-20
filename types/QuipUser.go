package types

type QuipUser struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Affinity          float64  `json:"affinity"`
	ProfilePictureURL *string  `json:"profile_picture_url"`
	ChatThreadID      *string  `json:"chat_thread_id"`
	DesktopFolderID   *string  `json:"desktop_folder_id"`
	ArchiveFolderID   *string  `json:"archive_folder_id"`
	StarredFolderID   *string  `json:"starred_folder_id"`
	PrivateFolderID   *string  `json:"private_folder_id"`
	GroupFolderIDs    []string `json:"group_folder_ids"`
	SharedFolderIDs   []string `json:"shared_folder_ids"`
	Disabled          *string  `json:"disabled"`
	CreatedUsec       *uint64  `json:"created_usec"`
}

func (qu QuipUser) Folders() []string {
	folders := make([]string, 0, 10)
	folders = tryAppend(
		folders,
		qu.ArchiveFolderID,
		qu.DesktopFolderID,
		qu.PrivateFolderID,
	)
	for _, folder := range qu.GroupFolderIDs {
		folders = append(folders, folder)
	}
	for _, folder := range qu.SharedFolderIDs {
		folders = append(folders, folder)
	}
	return folders
}

func tryAppend(list []string, items ...*string) []string {
	for _, item := range items {
		if item != nil {
			list = append(list, *item)
		}
	}

	return list
}
