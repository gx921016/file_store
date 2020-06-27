package meta

import (
	mydb "file_store/db"
	"sort"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)

}

//新增或更新元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

//更新数据到mysql
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFiniished(
		fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)

}

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	meta, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fileMeta := FileMeta{
		FileSize: meta.FileSize.Int64,
		FileName: meta.FileName.String,
		FileSha1: meta.FileHash,
		Location: meta.FileAddr.String,
		UploadAt: meta.FileUpdate.String,
	}

	return fileMeta, nil
}

func GetLastFileMetas(count int) []FileMeta {
	lastFileMetas := make([]FileMeta, 0)
	for _, v := range fileMetas {
		lastFileMetas = append(lastFileMetas, v)
	}
	sort.Sort(ByUploadTime(lastFileMetas))
	if count > len(lastFileMetas) {
		return lastFileMetas
	}
	return lastFileMetas[0:count]
}
