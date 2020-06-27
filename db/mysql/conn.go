package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlDb *sql.DB
var MysqlDbErr error

const (
	USER_NAME = "root"
	PASS_WORD = "gx921016"
	HOST      = "localhost"
	PORT      = "3306"
	DATABASE  = "filecloud"
	CHARSET   = "utf8"
)

func init() {
	dataBaseInit()
}

// 初始化链接
func dataBaseInit() {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", USER_NAME, PASS_WORD, HOST, PORT, DATABASE, CHARSET)

	// 打开连接失败
	mysqlDb, MysqlDbErr = sql.Open("mysql", dbDSN)
	//defer MysqlDb.Close();
	if MysqlDbErr != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + MysqlDbErr.Error())
	}

	// 最大连接数
	mysqlDb.SetMaxOpenConns(100)
	// 闲置连接数
	mysqlDb.SetMaxIdleConns(20)
	// 最大连接周期
	mysqlDb.SetConnMaxLifetime(100 * time.Second)

	if MysqlDbErr = mysqlDb.Ping(); nil != MysqlDbErr {
		panic("数据库链接失败: " + MysqlDbErr.Error())
	}

}

func DBConn() *sql.DB {
	return mysqlDb
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	// 获取记录列(名)
	columns, _ := rows.Columns()
	// 创建列值的slice (values)，并为每一列初始化一个指针
	// scanArgs用作rows.Scan中的传入参数
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	// record为每次迭代中存储行记录的临时变量
	record := make(map[string]interface{})
	// records为函数最终返回的数据(列表)
	records := make([]map[string]interface{}, 0)
	// 迭代行记录
	for rows.Next() {
		//每Scan一次，将一行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
