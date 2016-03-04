package uploader

type FileGroup struct {
	Id    int64  `db:"group_id" json:"id"`
	Name  string `db:"group_name" json:"name"`
	Count int64  `db:"cnt" json:"count"`
}

type FileGroupList struct {
	Uin    int64       `json:"biz_uin"`
	Groups []FileGroup `json:"file_group"`
}

type FileItem struct {
	FileId         string `json:"file_id"`
	Name           string `json:"name"`
	Type           int64  `json:"type"`
	Size           int64  `json:"size"`
	LastUpdateTime int64  `json:"update_time"`
	CdnUrl         string `json:"cdn_url"`
	ImageFormat    string `json:"image_format"`
}

type FileCount struct {
	Total  int64 `json:"total"`
	ImgCnt int64 `json:"img_cnt"`
}

type FilePageInfo struct {
	Type          int64         `json:"type"`
	FileCount     FileCount     `json:"file_cnt"`
	FileItems     []FileItem    `json:"file_item"`
	FileGroupList FileGroupList `json:"file_group_list"`
}

type UploadedFileInfo struct {
	FileId           string `json:"fileId"`
	FilePath         string `json:"filePath"`
	FileType         string `json:"fileType"`
	OriginalFileName string `json:"originalFileName"`
	FileSize         int64  `json:"fileSize"`
	CreateTime       int64  `json:"createTime"`
}
