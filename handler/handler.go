package handler

import (
	"encoding/json"
	"file_store/meta"
	"file_store/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHander(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		file, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(file))

	} else if r.Method == "POST" {
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMate := meta.FileMeta{
			FileName: header.Filename,
			Location: "/Users/gaoxiang/Downloads/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMate.Location)
		if err != nil {
			fmt.Printf("Failed to create file err:%s\n", err.Error())
		}
		defer newFile.Close()
		fileMate.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to copy file err:%s\n", err.Error())
		}

		newFile.Seek(0, 0)
		sha1 := util.FileSha1(newFile)
		fileMate.FileSha1 = sha1
		_ = meta.UpdateFileMetaDB(fileMate)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Successed")
}
func FileQueryMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	atoi, _ := strconv.Atoi(r.Form.Get("limit"))

	fileMeta := meta.GetLastFileMetas(atoi)
	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(marshal)

}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	open, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer open.Close()
	all, err := ioutil.ReadAll(open)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "appliication/octect-steam")
	w.Header().Set("content-disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	w.Write(all)

}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(marshal)
}
