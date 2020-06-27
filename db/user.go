package db

import (
	mydb "file_store/db/mysql"
	"fmt"
	"log"
)

// User : 用户表model
type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func UserSignup(username string, password string) bool {
	sqlStr := "insert ignore into tbl_user(`user_name`,`user_pwd`) values(?,?)"
	stmt, err := mydb.DBConn().Prepare(sqlStr)
	if err != nil {
		log.Println("Failed to insert err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		log.Println("Failed to insert to err:" + err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	log.Println(ret.RowsAffected())
	return false
}

//判断密码是否一致
func UserSinin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from  tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		log.Println(err.Error())
		return false
	} else if rows == nil {
		log.Println("username not found:" + username)
		return false
	}
	defer rows.Close()
	parseRows := mydb.ParseRows(rows)
	if len(parseRows) > 0 && string(parseRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare("replace into tbl_user_token(`user_name`,`user_token`) value(?,?)")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func GetUserToken(username string) (string, error) {

	stmt, err := mydb.DBConn().Prepare(
		"select user_token from tbl_user_token where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer stmt.Close()
	var token string
	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetUserInfo : 查询用户信息
func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
