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

	"golang.org/x/crypto/bcrypt"
)

func GetDogonData(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		request, err := http.Get("http://localhost:7070/logons")
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var logon []Logon
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &logon)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		var ranges = len(logon)
		fmt.Println("Sukses Ambil Data Login, Jumlah Data =", ranges)
		t, err := template.ParseFiles("website/templates/logon/logon.html")
		if err != nil {
			fmt.Println("Tidak bisa load template ")
		}
		t.Execute(w, logon)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func DeleteLogon(w http.ResponseWriter, r *http.Request) {
	newId := r.FormValue("userid")
	client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
	request, err := http.NewRequest("DELETE", "http://localhost:7070/logon/"+newId, nil)
	if err != nil {
		fmt.Println("Cannot call api delete")
		return
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Cannot do request")
		return
	}
	defer response.Body.Close()
	http.Redirect(w, r, "/logon/data", 301)
	fmt.Println("Succes Delete :", newId)
}

func GantiPasswordAdmin(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		if r.Method != "POST" {
			t, err := template.ParseFiles("website/templates/logon/logonubahpassword.html")
			if err != nil {
				fmt.Println("Tidak bisa load halaman")
			}
			t.Execute(w, nil)
		} else {
			response, err := http.Get("http://localhost:7070/logon/" + userName)
			if err != nil {
				fmt.Println("Tidak bisa mengambil data logon")
			}
			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Tidak bisa membaca data logon")
			}
			var logon Logon
			var jsonData = []byte(responseData)
			err = json.Unmarshal(jsonData, &logon)

			passwordlama := html.EscapeString(r.FormValue("passwordlama"))
			passwordbaru := html.EscapeString(r.FormValue("passwordbaru"))

			err = bcrypt.CompareHashAndPassword([]byte(logon.Password), []byte(passwordlama))
			if err != nil {
				fmt.Fprintf(w, "<script>alert('Password lama yang anda masukkan tidak sesuai !');window.location='/logon/ubahpwadmin';</script>")
			} else {
				var data = map[string]string{
					"password": passwordbaru,
				}
				var dataJSON, _ = json.Marshal(data)
				client := &http.Client{}
				request, err := http.NewRequest("PUT", "http://localhost:7070/logon/"+userName, bytes.NewBuffer(dataJSON))
				if err != nil {
					fmt.Fprintf(w, "<script>alert('Gagal ubah password !');window.location='/logon/ubahpwadmin';</script>")
					return
				} else {
					response, err := client.Do(request)
					if err != nil {
						fmt.Println("Cannot do request")
						return
					}
					fmt.Fprintf(w, "<script>alert('Password berhasil diubah !');window.location='/logon/ubahpwadmin';</script>")
					fmt.Println("Success ubah password pada username: ", userName)
					defer response.Body.Close()
				}
			}
			t, err := template.ParseFiles("website/templates/logon/logonubahpassword.html")
			if err != nil {
				fmt.Println("Tidak bisa load halaman")
			}
			t.Execute(w, nil)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
