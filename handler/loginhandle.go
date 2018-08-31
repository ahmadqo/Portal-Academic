package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/securecookie"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Id       int
	Username string
	PwdUser  string
}
type Session struct {
	Id        string
	CreatedAt time.Time
	User      Users
}
type Logon struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var sessions = map[string]Session{}
var users = map[string]Users{}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func GetUserName(r *http.Request) (userName string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, w http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

// ===========================================================
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		message := "Enter username and password to login!"
		retry := r.URL.Query().Get("retry")
		checkRetry, _ := strconv.ParseBool(retry)
		varmap := map[string]interface{}{
			"Message": message,
			"Status":  "",
		}
		if checkRetry == true {
			message = "Invalid Username or Password!"
			varmap["Message"] = message
			varmap["Status"] = "error"
		}
		t, _ := template.ParseFiles("website/templates/index.html")
		t.Execute(w, varmap)
		return
	} else {
		r.ParseForm()
		username := html.EscapeString(r.FormValue("username"))
		password := html.EscapeString(r.FormValue("password"))

		request, err := http.Get("http://localhost:7070/logon/" + username)
		if err != nil {
			fmt.Println("Tidak bisa mengambil data user")
			fmt.Fprintf(w, "<script>alert('Akun tidak bisa digunakan untuk login!');window.location='/';</script>")
		} else {
			responseData, err := ioutil.ReadAll(request.Body)
			var logon Logon
			var jsonData = []byte(responseData)
			err = json.Unmarshal(jsonData, &logon)
			if logon.Username == "" {
				fmt.Println("Tidak bisa mengambil data user")
				fmt.Fprintf(w, "<script>alert('Akun tidak bisa digunakan untuk login!');window.location='/';</script>")
			}
			// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dbPassword), bcrypt.DefaultCost)
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Println("passwordhas", string(hashedPassword))
			err = bcrypt.CompareHashAndPassword([]byte(logon.Password), []byte(password))
			if err != nil {
				fmt.Println("Tidak bisa compare password")
				fmt.Fprintf(w, "<script>alert('Password login salah!');window.location='/';</script>")
				return
			}

			log.Println(time.Now().Format(time.RFC850), "User Login Attempt by: ", username)
			fmt.Println("username :", username)
			fmt.Println("password :", logon.Password)
			setSession(username, w)
			if username == "admin" {
				http.Redirect(w, r, "/santri/index", 301)
			} else {
				http.Redirect(w, r, "/portal/index", 301)
			}
		}
	}
}

// ==========================================================

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	username := GetUserName(r)
	log.Println(time.Now().Format(time.RFC850), "User Logout Attempt by: ", username)
	clearSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "website/templates/index.html")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	jsonData := map[string]string{
		"username": username,
		"password": password,
	}
	jsonValue, _ := json.Marshal(jsonData)
	_, err := http.Post("http://localhost:7070/logon/register", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Fprintf(w, "<script>alert('Gagal Membuat Akun, akun sudah ada!');window.location='/';</script>")
	} else {
		fmt.Fprintf(w, "<script>alert('Akun Berhasil Dibuat!');window.location='/';</script>")
	}
}

// func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
// 	t, err := template.ParseFiles(filename)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	if err := t.Execute(w, data); err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// }

// func init() {
// 	t := time.NewTicker(600 * time.Second)
// 	go func() {
// 		for {
// 			for _, s := range sessions {
// 				if time.Now().Sub(s.CreatedAt) > (time.Second * 599) {
// 					delete(sessions, s.Id)
// 				}
// 			}
// 			<-t.C
// 		}
// 	}()
// }

// func ValidateUser(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		c, err := r.Cookie("session")
// 		if err != nil {
// 			http.Redirect(w, r, "/login", http.StatusSeeOther)
// 			fmt.Println("tidak bisa mengambil nama sessions")
// 			return
// 		}
// 		s, ok := sessions[c.Value]
// 		if !ok {
// 			http.Redirect(w, r, "/login", http.StatusSeeOther)
// 			fmt.Println("tidak bisa mengambil nama sessions 2")
// 			return
// 		}
// 		s.CreatedAt = time.Now()
// 		sessions[c.Value] = s

// 		u := s.User
// 		ctx := context.WithValue(r.Context(), "user", u)
// 		r = r.WithContext(ctx)

// 		h(w, r)
// 	}
// }
