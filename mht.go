package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var (
	change = make(map[string][]byte)
	//puser=make(map[int]string)
	//ppasswd=make(map[int]string)
	i, temp int = 0, 0
	userMap     = make(map[string]string)
)

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Raw        string
	Unparsed   []string
}

func init() {
	loadHtml("login", "/home/guaniu/code/src/http/login.html")
	loadHtml("home", "/home/guaniu/code/src/http/home.html")
	loadHtml("err", "/home/guaniu/code/src/http/err.html")
	loadHtml("reg", "/home/guaniu/code/src/http/reg.html")
	loadHtml("errtwo", "/home/guaniu/code/src/http/errtwo.html")
}
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Login)
	router.HandleFunc("/home", home)
	router.HandleFunc("/err", err)
	router.HandleFunc("/reg", reg)
	router.HandleFunc("/errtwo", errtwo)
	log.Fatal(http.ListenAndServe(":8080", router))
	/*
		http.HandleFunc("/", Login)
		http.HandleFunc("/home", home)
		http.HandleFunc("/err", err)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Print(err)
	} */
}
func Login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	passwd := r.FormValue("passwd")

	if len(user) == 0 || len(passwd) == 0 {
		fmt.Fprintf(w, "%s", change["login"])
	} else {

		pw, ok := userMap[user]
		if ok == true {
			if pw == passwd {

				cookieUser := http.Cookie{Name: "user", Value: user, Path: "/", MaxAge: 86400}
				cookiePasswd := http.Cookie{Name: "password", Value: passwd, Path: "/", MaxAge: 86400}
				http.SetCookie(w, &cookieUser)
				http.SetCookie(w, &cookiePasswd)
				http.Redirect(w, r, "/home", http.StatusFound)

			} else {
				http.Redirect(w, r, "/err", http.StatusFound)
			}
		} else {
			http.Redirect(w, r, "/err", http.StatusFound)
		}

		// for j:=0; j<=i ;j++ {
		// 	if user==puser[j]&& passwd==ppasswd[j]{
		// 		temp=1
		// 		break
		// 	}
		// }
		// if temp==1{
		// 	cookieUser := http.Cookie{Name:"user",Value:user,Path:"/",MaxAge:86400}
		// 	cookiePasswd := http.Cookie{Name:"password",Value:passwd,Path:"/",MaxAge:86400}
		// 	http.SetCookie(w,&cookieUser)
		// 	http.SetCookie(w,&cookiePasswd)
		// 	http.Redirect(w, r, "/home", http.StatusFound)
		// 	temp=0
		// }else{
		// 	http.Redirect(w, r, "/err", http.StatusFound)
		// }
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "%s", change["home"])
	cookieUser, err := r.Cookie("user")
	if err == nil {
		cookieUserValue := cookieUser.Value
		//将数据写入http连接中
		w.Write([]byte("用户名为：" + cookieUserValue))
	} else {
		w.Write([]byte("读取cookie信息出错：" + err.Error()))
	}
	cookiePasswd, err := r.Cookie("password")
	if err == nil {
		cookiePasswdValue := cookiePasswd.Value
		//将数据写入http连接中
		w.Write([]byte("密码为：" + cookiePasswdValue))

		// cookie := http.Cookie{Name: "testcookiename", Path: "/", MaxAge: -1}
		// http.SetCookie(w, &cookie)
	} else {
		w.Write([]byte("读取cookie信息出错：" + err.Error()))
	}
}
func err(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", change["err"])
}
func errtwo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", change["errtwo"])
}
func reg(w http.ResponseWriter, r *http.Request) {
	ruser := r.FormValue("ruser")
	rpasswd := r.FormValue("rpasswd")
	rpasswdtwo := r.FormValue("rpasswdtwo")
	if len(ruser) == 0 || len(rpasswd) == 0 || len(rpasswdtwo) == 0 {
		fmt.Fprintf(w, "%s", change["reg"])
	} else {
		if rpasswd != rpasswdtwo {
			http.Redirect(w, r, "/errtwo", http.StatusFound)
		} else {
			// puser[i]=ruser
			// ppasswd[i]=rpasswd
			// i++

			userMap[ruser] = rpasswd

			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}
func loadHtml(key, file_name string) {
	info, err := readFile(file_name)
	if err != nil {
		fmt.Print(err)
		return
	}
	change[key] = info
}
func readFile(file_name string) ([]byte, error) {
	fi, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}

