package api

import (
	"encoding/json"
	"net/http"
	"tools"
		"models"
	"time"
	"io/ioutil"
	"db"
	)

func Guid(writer http.ResponseWriter, request *http.Request) {
	newGUID := tools.NewUUID()

	value, err := json.Marshal(&newGUID)
	if err != nil {
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Write(value)
}

func UploadFile(writer http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(1024*1024*1024)
	file := request.MultipartForm.File["uploaded"][0]
	fileReader,err := file.Open()
	if err != nil {
		return
	}

	data,err := ioutil.ReadAll(fileReader)
	if err != nil {
		return
	}

	fileSize := len(data)
		fileId := tools.NewUUID().Content.(string)
		fileName := file.Filename
		fileType := file.Header.Get("Content-Type")
		result,err := tools.SaveFileToDisk(fileId,data)
		if err != nil || !result{
		return
	}

	currentFile := models.File{
		FileId:fileId,
		FileName:fileName,
		FileSize:float64(fileSize),
		FileType:fileType,
		Url: request.Host + "/static/" + fileId,
	}

	currentFile.FileCreatedAt = time.Now().String()

	_,err = db.CreateFileObject(currentFile)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	resp := models.Result{
		Date:time.Now().String(),
		Content:currentFile,
	}

	value,err := json.Marshal(&resp)
	if err != nil{
		return
	}

	writer.Header().Set("Content-Type","application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Write(value)
}
