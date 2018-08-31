package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Nilai struct {
	KodeNilai     string `json:"kodenilai"`
	Nis           string `json:"nis"`
	KodePelajaran string `json:"kodepelajaran"`
	Nilai         int    `json:"nilai"`
}

type Pelajaran struct {
	KodePelajaran string `json:"kodepelajaran"`
	NamaPelajaran string `json:"namapelajaran"`
}
type AmbilNilai struct {
	KodeNilai     string `json:"kodenilai"`
	Nis           string `json:"nis"`
	Nama          string `json:"nama"`
	Komplek       string `json:"komplek"`
	KodePelajaran string `json:"kodepelajaran"`
	NamaPelajaran string `json:"namapelajaran"`
	Nilai         int    `json:"nilai"`
}

func HandleNilaiIndex(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		t, err := template.ParseFiles("website/templates/nilai/nilai.html")
		if err != nil {
			fmt.Println("error")
		}
		t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func HandeInputDataNilai(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		if r.Method != "POST" {
			http.ServeFile(w, r, "website/templates/nilai/nilaiinput.html")
			return
		}
		getnis := r.FormValue("nis")
		request, err := http.Get("http://localhost:8077/api/munawwir/santri/" + getnis)
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, _ := ioutil.ReadAll(request.Body)
		if len(responseData) == 0 {
			fmt.Fprintf(w, "<script>alert('Nomor Induk Santri Tidak Ditemukan!');window.location='/nilai/inputData-nilai';</script>")
		} else {
			kodenilai := r.FormValue("kodenilai")
			nis := r.FormValue("nis")
			kodepelajaran := r.FormValue("kodepelajaran")
			nilai := r.FormValue("nilai")
			var numnilai, _ = strconv.Atoi(nilai)

			jsonData := Nilai{kodenilai, nis, kodepelajaran, numnilai}
			jsonValue, _ := json.Marshal(jsonData)
			request, err = http.Post("http://localhost:8055/api/santri/nilai/insert", "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				fmt.Println("The http request failed")
			} else {
				data, _ := ioutil.ReadAll(request.Body)
				fmt.Println(string(data))
				http.Redirect(w, r, "/nilai/inputData-nilai", 301)
			}
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func GetTranskipNilaiSantri(w http.ResponseWriter, r *http.Request) {
	getNis := r.FormValue("nis")
	request, err := http.Get("http://localhost:8077/api/munawwir/santri/" + getNis)
	if err != nil {
		fmt.Println(err.Error())
	}
	responseData2, _ := ioutil.ReadAll(request.Body)
	if len(responseData2) == 0 {
		fmt.Fprintf(w, "<script>alert('Nomor Induk Santri Tidak Ditemukan!');window.location='/nilai/index';</script>")
	} else {
		request, err := http.Get("http://localhost:8055/api/santri/nilai/" + getNis)
		if err != nil {
			fmt.Println("Tidak Bisa Memanggil Api mengambil data nilai")
			return
		}
		responseData, _ := ioutil.ReadAll(request.Body)
		if len(responseData) == 0 {
			fmt.Fprintf(w, "<script>alert('Nomor Induk Santri Tidak Ditemukan!');window.location='/nilai/inputData-nilai';</script>")
		}
		var transkipnilai []AmbilNilai
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &transkipnilai)
		if err != nil {
			fmt.Println("Tidak bisa Unmarshal Data")
		}
		t, err := template.ParseFiles("website/templates/nilai/nilaitampil.html")
		if err != nil {
			fmt.Println("Tidak Bisa load halaman tampil transkip nilai")
		}
		t.Execute(w, transkipnilai)
	}
}

func UpdateNilaiSantri(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		getkode := r.FormValue("kode")
		request, err := http.Get("http://localhost:8055/api/santri/nilai/bykode/" + getkode)
		if err != nil {
			fmt.Println("Tidak Bisa memanggil API nilai by kode")
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println("Tidak bisa membaca body data nilai ")
		}
		var nilai Nilai
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &nilai)
		if err != nil {
			fmt.Println("Tidak bisa decode data nilai")
		}
		t, err := template.ParseFiles("website/templates/nilai/nilaiupdate.html")
		if err != nil {
			fmt.Println("Tidak bisa load template nilai update")
		}
		t.Execute(w, nilai)
		return
	} else {
		var kodenilai = r.FormValue("kodenilai")
		nis := r.FormValue("nis")
		kodepelajaran := r.FormValue("kodepelajaran")
		nilai := r.FormValue("nilai")
		var numnilai, _ = strconv.Atoi(nilai)

		jsonData := Nilai{kodenilai, nis, kodepelajaran, numnilai}
		jsonValue, _ := json.Marshal(jsonData)
		client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
		request, err := http.NewRequest("PUT", "http://localhost:8055/api/santri/nilai/update/"+kodenilai, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Cannot call api Updte")
			return
		}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Cannot do request")
			return
		} else {
			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error Tidak bisa Membaca data", err)
			} else {
				fmt.Println("Sukses Update Data Nilai Dengan Kode :", kodenilai)
				fmt.Println(string(data))
				fmt.Println()
				defer response.Body.Close()
				// ======================
				// fmt.Fprintf(w, "<script>alert('Update Data Nilai Sukses!');</script>")
				request, err := http.Get("http://localhost:8055/api/santri/nilai/" + nis)
				if err != nil {
					fmt.Println("Tidak Bisa Memanggil Api mengambil data nilai")
					return
				}
				responseData, _ := ioutil.ReadAll(request.Body)
				var transkipnilai []AmbilNilai
				var jsonData = []byte(responseData)
				err = json.Unmarshal(jsonData, &transkipnilai)
				if err != nil {
					fmt.Println("Tidak bisa Unmarshal Data")
				}
				t, err := template.ParseFiles("website/templates/nilai/nilaitampil.html")
				if err != nil {
					fmt.Println("Tidak Bisa load halaman tampil transkip nilai")
				}
				t.Execute(w, transkipnilai)
			}
		}
	}
}

func DeleteNilaiByKode(w http.ResponseWriter, r *http.Request) {
	getNis := r.FormValue("nis")
	getKode := r.FormValue("kode")
	client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
	request, err := http.NewRequest("DELETE", "http://localhost:8055/api/santri/nilai/delete/"+getKode, nil)
	if err != nil {
		fmt.Println("Tidak bisa memanggi method Delete Nilai")
		return
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Tidak bisa melakukan Do Request")
		return
	}
	defer response.Body.Close()
	fmt.Println("Sukses Delete Nilai Kode :", getKode)
	// ========================
	request2, err := http.Get("http://localhost:8055/api/santri/nilai/" + getNis)
	if err != nil {
		fmt.Println("Tidak Bisa Memanggil Api mengambil data nilai")
		return
	}
	responseData2, _ := ioutil.ReadAll(request2.Body)
	var transkipnilai []AmbilNilai
	var jsonData = []byte(responseData2)
	err = json.Unmarshal(jsonData, &transkipnilai)
	if err != nil {
		fmt.Println("Tidak bisa Unmarshal Data")
	}
	t, err := template.ParseFiles("website/templates/nilai/nilaitampil.html")
	if err != nil {
		fmt.Println("Tidak Bisa load halaman tampil transkip nilai")
	}
	t.Execute(w, transkipnilai)
}
