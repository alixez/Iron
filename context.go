package Iron

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/alixez/Iron/utils"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type (
	// APIResponse ...
	APIResponse struct {
		Code    int         `json:"Code" xml:"Code"`
		SubCode string      `json:"SubCode" xml:"SubCode"`
		Message string      `json:"Message" xml:"Message"`
		Data    interface{} `json:"Data" xml:"Data"`
	}

	// Context ...
	Context struct {
		echo.Context
		services    map[string]ServiceInterface
		apiResponse *APIResponse
		dbHelper    map[string]interface{}
		Config      *Env
	}

	// File ...
	File struct {
		Filename     string
		Path         string
		AbstructPath string
		Host         string
		Extension    string
	}
)

// SaveFilesToStorage 将文件保存到存储空间（多文件）
func (this *Context) SaveFilesToStorage(fields string, subpath string) ([]*File, error) {
	var fileList []*File
	form, err := this.MultipartForm()
	if err != nil {
		return nil, err
	}

	files := form.File[fields]

	for _, file := range files {
		if fileModel, err := this.executeUploadedFile(file, subpath); err == nil {
			fileList = append(fileList, fileModel)
		} else {
			return nil, err
		}

	}

	return fileList, nil
}

// SaveFileToStorage 将文件保存到存储空间(单文件)
func (ctx *Context) SaveFileToStorage(fields string, subpath string) (*File, error) {
	var fileModel *File

	file, err := ctx.FormFile(fields)
	if err != nil {
		return nil, err
	}
	fileModel, err = ctx.executeUploadedFile(file, subpath)
	if err != nil {
		return nil, err
	}
	return fileModel, nil
}

// ExecuteUploadedFile 处理已经上传的文件
func (ctx *Context) executeUploadedFile(file *multipart.FileHeader, subpath string) (*File, error) {
	config := ctx.Config
	storageInterface := config.Get("storage").(map[interface{}]interface{})
	storage := make(map[string]string)
	for k, v := range storageInterface {
		storage[k.(string)] = v.(string)
	}
	rootPath := storage["root"]
	tumbnailPath := filepath.Join(rootPath, storage["tumbnail"])
	orignailPath := filepath.Join(rootPath, storage["orignail"])
	mimeType := file.Header["Content-Type"][0]
	filename := uuid.NewV1().String() + "." + strings.Split(mimeType, "/")[1]
	dstPath := filepath.Join(orignailPath, subpath)

	if !utils.IsDirExist(rootPath) {
		os.Mkdir(rootPath, 0777)
	}
	if !utils.IsDirExist(tumbnailPath) {
		os.Mkdir(tumbnailPath, 0777)
	}
	if !utils.IsDirExist(orignailPath) {
		os.Mkdir(orignailPath, 0777)
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	if !utils.IsDirExist(dstPath) {

		os.MkdirAll(dstPath, 0777)
	}
	dst, err := os.Create(filepath.Join(dstPath, filename))
	if err != nil {
		return nil, err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(filepath.Join(dstPath, filename))
	if err != nil {
		return nil, err
	}

	fileModel := &File{
		Filename:     filename,
		Path:         filepath.Join(dstPath, filename),
		AbstructPath: absPath,
		Host:         storage["host"],
		Extension:    strings.Split(mimeType, "/")[1],
	}

	return fileModel, nil
}

func (this *Context) AddDBHelper(name string, value interface{}) {
	this.dbHelper[name] = value
}

func (this *Context) GetDB(name string) interface{} {
	return this.dbHelper[name]
}

func (this *Context) SetServices(services map[string]ServiceInterface) {

	this.services = services
}

func (this *Context) GetService(name string) ServiceInterface {
	service, ok := this.services[name]
	if ok == false {
		return nil
	}
	service.Init(this)
	return service
}
