package controllers

import (
	"encoding/json"
	"io"
	"os"

	"github.com/MarxGo/h5upload/util"

	"github.com/astaxie/beego"
)

var (
	// receiveFileMap : fileblockSequence : fileBlockName
	receiveFileMap map[int64]string
)

type UploadController struct {
	beego.Controller
}

func init() {
	receiveFileMap = make(map[int64]string)
}

// ToUpload is
func (upload *UploadController) ToUpload() {
	upload.TplName = "upload/upload.html"
}

// CheckFileExist is
func (upload *UploadController) CheckFileExist() {
	receiveFileMap = make(map[int64]string)
	result := make(map[string]interface{})

	reqMap := make(map[string]interface{})
	json.Unmarshal(upload.Ctx.Input.RequestBody, &reqMap)

	filePath := "upload/" + reqMap["id"].(string) + "/" + reqMap["id"].(string)
	if util.IsPathExist(filePath) {
		isExist, err := util.CheckFileMd5(filePath, reqMap["id"].(string))
		result["isExist"] = isExist
		if err != nil {
			beego.Error(err)
		}
	} else {
		result["isExist"] = false
	}

	upload.Data["json"] = result
	upload.ServeJSON()
}

// CheckFileBlockExist is
func (upload *UploadController) CheckFileBlockExist() {
	result := make(map[string]interface{})

	reqMap := make(map[string]interface{})
	json.Unmarshal(upload.Ctx.Input.RequestBody, &reqMap)

	blockFilePath := "upload/" + reqMap["id"].(string) + "/" + reqMap["blockId"].(string)
	if util.IsPathExist(blockFilePath) {
		isExist, err := util.CheckFileMd5(blockFilePath, reqMap["blockId"].(string))
		result["isExist"] = isExist
		if err != nil {
			beego.Error(err)
		}
	} else {
		result["isExist"] = false
	}

	upload.Data["json"] = result
	upload.ServeJSON()
}

// ReceiveFile is
func (upload *UploadController) ReceiveFile() {
	result := make(map[string]string)
	fileId := upload.GetString("fileId")
	blockId := upload.GetString("blockId")
	blockCount, _ := upload.GetInt64("currentBlock")
	blockSum, _ := upload.GetInt64("sumBlock")

	dirPath := "upload/" + fileId
	util.CompleteDirPath(dirPath)

	err := upload.SaveToFile("data", dirPath+"/"+blockId)
	if err != nil {
		beego.Error(err)
		result["msg"] = "error when save file:" + err.Error()
		upload.Data["json"] = result
		upload.ServeJSON()
		upload.StopRun()
	}

	receiveFileMap[blockCount] = blockId

	// if receive all block of current file,merge it
	if int64(len(receiveFileMap)) == blockSum {
		mergeFile(fileId, blockSum)
	}

	result["msg"] = "success"
	upload.Data["json"] = result
	upload.ServeJSON()
}

// GetBlockSizeAndWorkerNum is
func (upload *UploadController) GetBlockSizeAndWorkerNum() {
	result := make(map[string]interface{})

	reqMap := make(map[string]interface{})
	json.Unmarshal(upload.Ctx.Input.RequestBody, &reqMap)

	result["income"] = reqMap
	result["workerNum"] = 3
	result["fileSize"] = 128

	upload.Data["json"] = result
	upload.ServeJSON()
}

// Empty is
func (upload *UploadController) Empty() {
	result := make(map[string]string)
	result["result"] = "success"
	fileId := upload.GetString("fileId")
	blockId := upload.GetString("blockId")
	blockCount, _ := upload.GetInt64("currentBlock")
	blockSum, _ := upload.GetInt64("sumBlock")

	receiveFileMap[blockCount] = blockId
	if int64(len(receiveFileMap)) == blockSum {
		mergeFile(fileId, blockSum)
	}

	upload.Data["json"] = result
	upload.ServeJSON()
}

func mergeFile(fileId string, blockSum int64) {
	defer func() {
		receiveFileMap = make(map[int64]string)
	}()
	filePath := "upload/" + fileId + "/" + fileId
	finalFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		beego.Error(err)
		return
	}
	defer finalFile.Close()

	for i := int64(0); i < blockSum; i++ {
		buf := make([]byte, 100)
		blockFilePath := "upload/" + fileId + "/" + receiveFileMap[i]

		blockFile, err := os.Open(blockFilePath)
		if err != nil {
			beego.Error(err)
			return
		}
		defer blockFile.Close()

		for {
			n, err := blockFile.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if 0 == n {
				break
			}
			finalFile.Write(buf[:n])
			buf = make([]byte, 100)
		}
	}
}
