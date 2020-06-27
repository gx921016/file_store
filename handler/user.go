package handler

import (
	dblayer "file_store/db"
	"file_store/util"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	pwd_salt = "^#610"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Printf("MethodGet")
		file, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(file)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}
	enc_passw := util.Sha1([]byte(passwd + pwd_salt))
	suc := dblayer.UserSignup(username, enc_passw)
	if suc {
		w.Write([]byte("SUCCESS"))
		log.Printf("SUCCESS")
	} else {
		log.Printf("FAILED")
		w.Write([]byte("FAILED"))
	}
}

//登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	//1.效验用户名和密码
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPassw := util.Sha1([]byte(password + pwd_salt))
	pwdChecked := dblayer.UserSinin(username, encPassw)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
	}
	//2.生成访问凭证
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		return
	}
	log.Println("登录成功后重定向到首页" + "http://" + r.Host + "/static/view/home.html")
	//3.登录成功后重定向到首页
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	isValid := isTokenValid(token, username)
	if !isValid {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}
func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

func isTokenValid(token string, username string) bool {
	if len(token) != 40 {
		return false
	}

	// TODO: 判断token的时效性，是否过期
	// example，假设token的有效期为1天   (根据同学们反馈完善, 相对于视频有所更新)
	tokenTS := token[:8]
	if util.Hex2Dec(tokenTS) < time.Now().Unix()-86400 {
		return false
	}

	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	// example, IsTokenValid方法增加传入参数username
	userToken, _ := dblayer.GetUserToken(username)
	if userToken != token {
		return false
	}

	return true
}
