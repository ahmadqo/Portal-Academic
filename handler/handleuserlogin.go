package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	*sql.DB
}

func NewConnect() (DB, error) {
	connect := "postgres://root@127.0.0.1:26257/santri_db?sslmode=disable"
	db, err := sql.Open("postgres", connect)
	if err != nil {
		log.Fatal("Cannot Connect DB ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Tidak terhubung ke database", err)
	}
	return DB{db}, err
}

func OpenIndexSantri(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		t, err := template.ParseFiles("website/templateuser/dashbord.html")
		if err != nil {
			fmt.Println("error")
		}
		t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func UserGetDataKegiatan(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		request, err := http.Get("http://localhost:8090/api/santri/kegiatan/datas")
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var kegiatan []Kegiatan
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &kegiatan)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		var ranges = len(kegiatan)
		fmt.Println("Sukses Ambil Data Kegaitan, Data =", ranges)
		t, err := template.ParseFiles("website/templateuser/user_kegiatan.html")
		if err != nil {
			fmt.Println("Tidak bisa load template")
		}
		t.Execute(w, kegiatan)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func UserGetDataKegiatanHarian(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		request, err := http.Get("http://localhost:8090/api/kegiatan/harian/santris")
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var kegharian []KegiatanHarian
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &kegharian)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		var ranges = len(kegharian)
		fmt.Println("Sukses Ambil Data Kegaitan, Data =", ranges)
		t, err := template.ParseFiles("website/templateuser/user_kegiatan_harian.html")
		if err != nil {
			fmt.Println("Tidak bisa load template")
		}
		t.Execute(w, kegharian)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func UserGetDataKajian(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		request, err := http.Get("http://localhost:8088/api/santri/kajian/datas")
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var kajians []Kajian
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &kajians)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		var ranges = len(kajians)
		fmt.Println("Sukses Ambil Data Kajian, Jumlah Data =", ranges)
		t, err := template.ParseFiles("website/templateuser/user_kajian.html")
		if err != nil {
			fmt.Println("Tidak bisa load template")
		}
		t.Execute(w, kajians)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func UserGetDataSantriInLogin(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		request, err := http.Get("http://localhost:8077/api/munawwir/santri/" + userName)
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var santri Santri
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &santri)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		t, err := template.ParseFiles("website/templateuser/user_santri.html")
		if err != nil {
			fmt.Println("Tidak bisa load template Halaman Santri detil")
		}
		t.Execute(w, santri)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func UserGetTranskipNilaiSantri(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		request, err := http.Get("http://localhost:8077/api/munawwir/santri/" + userName)
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData2, _ := ioutil.ReadAll(request.Body)
		if len(responseData2) == 0 {
			fmt.Fprintf(w, "<script>alert('Nomor Induk Santri Tidak Ditemukan!');window.location='/portal/index';</script>")
		} else {
			request, err := http.Get("http://localhost:8055/api/santri/nilai/" + userName)
			if err != nil {
				fmt.Println("Tidak Bisa Memanggil Api mengambil data nilai")
				return
			}
			responseData, _ := ioutil.ReadAll(request.Body)
			// if len(responseData) == 0 {
			// 	fmt.Fprintf(w, "<script>alert('Nomor Induk Santri Tidak Ditemukan!');window.location='/nilai/inputData-nilai';</script>")
			// }
			var transkipnilai []AmbilNilai
			var jsonData = []byte(responseData)
			err = json.Unmarshal(jsonData, &transkipnilai)
			if err != nil {
				fmt.Println("Tidak bisa Unmarshal Data")
			}
			t, err := template.ParseFiles("website/templateuser/user_nilai.html")
			if err != nil {
				fmt.Println("Tidak Bisa load halaman")
			}
			t.Execute(w, transkipnilai)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

}

func UserJurnal(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		request, err := http.Get("http://localhost:8098/jurnal/getdata")
		if err != nil {
			fmt.Println("Tidak bisa memanggil API")
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var jurs []Jurnal
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &jurs)
		t, err := template.ParseFiles("website/templateuser/user_jurnal.html")
		if err != nil {
			fmt.Println("Tidak Bisa load halaman")
		}
		t.Execute(w, jurs)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func UserGantiPassword(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		if r.Method != "POST" {
			t, err := template.ParseFiles("website/templateuser/user_password.html")
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
				fmt.Fprintf(w, "<script>alert('Password lama yang anda masukkan tidak sesuai !');window.location='/portal/santri/ubahpassword';</script>")
			} else {
				var data = map[string]string{
					"password": passwordbaru,
				}
				var dataJSON, _ = json.Marshal(data)
				client := &http.Client{}
				request, err := http.NewRequest("PUT", "http://localhost:7070/logon/"+userName, bytes.NewBuffer(dataJSON))
				if err != nil {
					fmt.Fprintf(w, "<script>alert('Gagal ubah password !');window.location='/portal/santri/ubahpassword';</script>")
					return
				} else {
					response, err := client.Do(request)
					if err != nil {
						fmt.Println("Cannot do request")
						return
					}
					fmt.Fprintf(w, "<script>alert('Password berhasil diubah !');window.location='/portal/santri/ubahpassword';</script>")
					fmt.Println("Success ubah password pada username: ", userName)
					defer response.Body.Close()
				}
			}
			t, err := template.ParseFiles("website/templateuser/user_password.html")
			if err != nil {
				fmt.Println("Tidak bisa load halaman")
			}
			t.Execute(w, nil)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
