package uploader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
	"github.com/swanwish/go-common/web"
)

type UploadHandler struct {
	FileUploadPath      string
	FuncSaveFileInfo    func(web.HandlerContext, UploadedFileInfo) (string, error)
	FuncGetFileInfoById func(fileId string) (UploadedFileInfo, error)
	FuncDeleteFileByIds func(fileId string) error
	FuncGetPageInfo     func(web.HandlerContext) (FilePageInfo, error)
	FuncModifyFile      func(web.HandlerContext) error
	FuncPostFilePage    func(web.HandlerContext) error
}

func (h UploadHandler) UploadFiles(ctx web.HandlerContext) {
	fileInfo, err := handleUploadFile(ctx, "file", h.FileUploadPath)
	if err != nil {
		logs.Errorf("Failed to handle upload file")
		ctx.ReplyInvalidParameterError()
		return
	}
	if h.FuncSaveFileInfo != nil {
		fileInfo.FileId, err = h.FuncSaveFileInfo(ctx, fileInfo)
		if err != nil {
			logs.Errorf("Failed to save uploaded file information, the error is %v", err)
			ctx.ReplyInternalError()
			return
		}
	}
	ctx.ReplyJsonData(fileInfo)
}

func (h UploadHandler) DownloadFile(ctx web.HandlerContext) {
	if h.FuncGetFileInfoById == nil {
		logs.Errorf("The function FuncGetFileInfoById is not specified.")
		ctx.ReplyInternalError()
	}
	fileId := ctx.Var("fileId")
	fileInfo, err := h.FuncGetFileInfoById(fileId)
	if err != nil {
		ctx.ReplyInvalidParameterError()
		return
	}
	serverFilePath(ctx, fileInfo.FilePath)
}

func (h UploadHandler) GetFile(ctx web.HandlerContext) {
	folder := ctx.Var("folder")
	fileName := ctx.Var("fileName")
	filePath := fmt.Sprintf("%s/%s", folder, fileName)
	serverFilePath(ctx, filePath)
}

func serverFilePath(ctx web.HandlerContext, filePath string) {
	ctx.ServeFile(filePath)
}

func (h UploadHandler) DeleteUploadedFile(ctx web.HandlerContext) {
	if h.FuncDeleteFileByIds == nil {
		logs.Error("The delete file by id function is not specified.")
		ctx.ReplyInternalError()
		return
	}
	fileIds := ctx.Var("fileIds")
	err := h.FuncDeleteFileByIds(fileIds)
	if err != nil {
		ctx.ReplyInternalError()
		return
	}
	ctx.ReplyJsonData(nil)
}

func (h UploadHandler) GetFilePage(ctx web.HandlerContext) {
	if h.FuncGetPageInfo == nil {
		logs.Error("The function FuncGetPageInfo is not specified.")
		ctx.ReplyInternalError()
		return
	}
	pageInfo, err := h.FuncGetPageInfo(ctx)
	if err != nil {
		ctx.ReplyInternalError()
		return
	}
	ctx.ReplyJsonData(pageInfo)
	return
}

func (h UploadHandler) PostFilePage(ctx web.HandlerContext) {
	if h.FuncPostFilePage == nil {
		logs.Errorf("The FuncPostFilePage is not specified.")
		ctx.ReplyInternalError()
		return
	}
	err := h.FuncPostFilePage(ctx)
	if err != nil {
		logs.Errorf("Failed to execut FuncPostFilePage, the error is %v", err)
		ctx.ReplyInternalError()
		return
	}
	ctx.ReplyJsonData(nil)
}

func (h UploadHandler) ModifyFile(ctx web.HandlerContext) {
	if h.FuncModifyFile == nil {
		logs.Errorf("The function FuncModifyFile is not specified")
		ctx.ReplyInternalError()
		return
	}
	err := h.FuncModifyFile(ctx)
	if err != nil {
		ctx.ReplyInternalError()
		return
	}
	ctx.ReplyJsonData(nil)
}

func handleUploadFile(ctx web.HandlerContext, formFileName, fileUploadPath string) (UploadedFileInfo, error) {
	fileInfo := UploadedFileInfo{}
	ctx.R.ParseMultipartForm(32 << 20)
	inputFile, handler, err := ctx.R.FormFile(formFileName)
	if err != nil {
		logs.Errorf("Failed to parse form file, the error is %v", err)
		return fileInfo, err
	}
	defer inputFile.Close()
	originalFileName := handler.Filename
	fileInfo.OriginalFileName = originalFileName
	lastDotIndex := strings.LastIndex(originalFileName, ".")
	destFileName, err := utils.RandomKey()
	if err != nil {
		logs.Errorf("Failed to generate random file name.")
		return fileInfo, err
	}

	fileType := FILE_TYPE_FILES
	if lastDotIndex > 0 {
		extension := originalFileName[lastDotIndex+1:]
		fileType = GetFileType(extension)
		destFileName = fmt.Sprintf("%s.%s", destFileName, extension)
	}

	fileInfo.FileType = fileType

	now := time.Now()
	fileDir := filepath.Join(fileType, fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()))
	fileDest := filepath.Join(fileUploadPath, fileDir)

	if !utils.FileExists(fileDest) {
		err = os.MkdirAll(fileDest, 0777)
		if err != nil {
			logs.Errorf("Failed to create upload path, the error is %v", err)
			return fileInfo, err
		}
	}

	destFilePath := filepath.Join(fileDest, destFileName)
	fileInfo.FilePath = filepath.Join(fileDir, destFileName)
	f, err := os.OpenFile(destFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logs.Errorf("Failed to open file, the error is %v", err)
		return fileInfo, err
	}
	defer f.Close()
	written, err := io.Copy(f, inputFile)
	if err != nil {
		logs.Errorf("Failed to copy file, the error is %v", err)
		return fileInfo, err
	}
	fileInfo.FileSize = int64(written)
	//if settings.OssClient != nil {
	//	logs.Debugf("Will create object %s", fmt.Sprintf("/%s/%s", common.OSSImageBucket, destFilePath))
	//	err = common.OssClient.CreateObject(fmt.Sprintf("/%s/%s", common.OSSImageBucket, destFilePath), destFilePath)
	//	if err != nil {
	//		logs.Errorf("Failed to upload file to oss, the erros is %v", err)
	//		return fileInfo, err
	//	}
	//	logs.Debugf("After create object")
	//}
	return fileInfo, nil
}
