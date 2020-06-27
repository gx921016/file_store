package db

import (
	"database/sql"
	mydb "file_store/db/mysql"
	"fmt"
)

func OnFileUploadFiniished(filehash string, filename string,
	filesize int64, fileaddr string) bool {
	//insert ignore into会忽略已有数据
	//sqlStr:="insert ignore into tbl_file(file_sha1,file_name,file_size,file_addr,status) value(?,?,?,?,1)"
	prepare, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`,`status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer prepare.Close()
	exec, err := prepare.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if rf, err := exec.RowsAffected(); nil == err {
		if rf <= 0 {

		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash   string
	FileName   sql.NullString
	FileSize   sql.NullInt64
	FileAddr   sql.NullString
	FileUpdate sql.NullString
}

func GetFileMeta(filehash string) (*TableFile, error) {
	prepare, err := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,file_addr,update_at from tbl_file " +
			"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer prepare.Close()
	tfile := TableFile{}
	err = prepare.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr, &tfile.FileUpdate)
	if err != nil {
		return nil, err
	}
	return &tfile, err
}
