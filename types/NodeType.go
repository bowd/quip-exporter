package types

type NodeType = string

var NodeTypes = struct {
	CurrentUser       NodeType
	Blob              NodeType
	User              NodeType
	Folder            NodeType
	Thread            NodeType
	ThreadHTML        NodeType
	ThreadComments    NodeType
	ThreadDocument    NodeType
	ThreadSlides      NodeType
	ThreadSpreadsheet NodeType
	Archive           NodeType
	UserPicture       NodeType
	ThreadIndex       NodeType
	FolderIndex       NodeType
}{
	CurrentUser:       "current-user",
	Blob:              "blob",
	User:              "user",
	Folder:            "folder",
	Thread:            "thread",
	ThreadHTML:        "thread-html",
	ThreadComments:    "thread-comments",
	ThreadDocument:    "thread-document",
	ThreadSlides:      "thread-slides",
	ThreadSpreadsheet: "thread-spreadsheet",
	Archive:           "archive",
	UserPicture:       "user-picture",
	ThreadIndex:       "thread-index",
	FolderIndex:       "folder-index",
}
