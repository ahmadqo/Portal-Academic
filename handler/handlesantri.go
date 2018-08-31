package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Santri struct {
	Id           int    `json:"id"`
	Nis          string `json:"nis"`
	Nama         string `json:"nama"`
	JenisKelamin string `json:"jk"`
	TempatLahir  string `json:"tempatlahir"`
	TanggalLahir string `json:"tgllahir"`
	Alamat       string `json:"alamat"`
	TanggalMasuk string `json:"tglmasuk"`
	Komplek      string `json:"komplek"`
	NamaAyah     string `json:"nmayah"`
	NamaIbu      string `json:"nmibu"`
}

func HalamanSantri(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		t, err := template.ParseFiles("website/templates/santri/santri.html")
		if err != nil {
			fmt.Println("error")
		}
		t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func InputDataSantri(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		if r.Method != "POST" {
			http.ServeFile(w, r, "website/templates/santri/santriinput.html")
			return
		}
		var buffer bytes.Buffer
		var komplek string
		var nis string
		nodaftar := r.FormValue("nodaftar")
		namasantri := r.FormValue("nama")
		jk := r.FormValue("jk")
		tempatlahir := r.FormValue("tempatlahir")
		tgllahir := r.FormValue("tgllahir")
		alamat := r.FormValue("alamat")
		tglmasuk := r.FormValue("tglmasuk")
		komplekname := r.FormValue("komplek")
		namaayah := r.FormValue("namaayah")
		namaibu := r.FormValue("namaibu")

		substring := tglmasuk[2:len(tglmasuk)]
		inputFmt := substring[:len(substring)-6]

		buffer.WriteString(komplekname)
		buffer.WriteString(inputFmt)
		buffer.WriteString(nodaftar)
		nis = (buffer.String())

		switch komplekname {
		case "110":
			komplek = "Komplek AB"
		case "701":
			komplek = "Komplek Gipa"
		case "081":
			komplek = "Huffadz 1"
		case "082":
			komplek = "Huffadz 2"
		case "901":
			komplek = "Komplek IJ"
		case "012":
			komplek = "Komplek L Putra"
		case "013":
			komplek = "Komplek L Putri"
		case "015":
			komplek = "Komplek M"
		case "231":
			komplek = "Nurussalam Putra"
		case "232":
			komplek = "Nurussalam Putri"
		case "601":
			komplek = "Padang Jagad"
		case "005":
			komplek = "Komplek Q"
		case "061":
			komplek = "Komplek R1"
		case "062":
			komplek = "Komplek R2"
		case "056":
			komplek = "Ribatul Qur'an"
		default:
			komplek = ""
		}
		jsonData := map[string]string{
			"nis":         nis,
			"nama":        namasantri,
			"jk":          jk,
			"tempatlahir": tempatlahir,
			"tgllahir":    tgllahir,
			"alamat":      alamat,
			"tglmasuk":    tglmasuk,
			"komplek":     komplek,
			"nmayah":      namaayah,
			"nmibu":       namaibu}

		jsonValue, _ := json.Marshal(jsonData)
		request, err := http.Post("http://localhost:8077/api/munawwir/santri/data", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("The http request failed")
		} else {
			data, _ := ioutil.ReadAll(request.Body)
			fmt.Println(string(data))
		}
		http.Redirect(w, r, "/santri/dataAll-daftar", 301)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func GetAllDataSantri(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		request, err := http.Get("http://localhost:8077/api/munawwir/santri/datas")
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var santri []Santri
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &santri)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		var ranges = len(santri)
		fmt.Println("Sukses Ambil Data Santri, Data =", ranges)
		t, err := template.ParseFiles("website/templates/santri/santriall.html")
		if err != nil {
			fmt.Println("Tidak bisa load template Daftar Santri")
		}
		t.Execute(w, santri)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func GetDataSantriByNis(w http.ResponseWriter, r *http.Request) {
	var getNis = r.FormValue("nis")
	request, err := http.Get("http://localhost:8077/api/munawwir/santri/" + getNis)
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
	t, err := template.ParseFiles("website/templates/santri/santribynis.html")
	if err != nil {
		fmt.Println("Tidak bisa load template Halaman Santri detil")
	}
	t.Execute(w, santri)
}

func UpdateDataSantri(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		var getNis = r.FormValue("nis")
		var getId = r.FormValue("id")
		request, err := http.Get("http://localhost:8077/api/munawwir/santri/" + getNis)
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
		if getId == "0" {
			http.Redirect(w, r, "/santri/dataAll-daftar", 301)
		} else {
			t, err := template.ParseFiles("website/templates/santri/santriupdate.html")
			if err != nil {
				fmt.Println("Tidak Dapat me-Load Halaman Update Santri")
			}
			t.Execute(w, santri)
			return
		}
	} else {
		newId := r.FormValue("id")
		nis := r.FormValue("nis")
		namasantri := r.FormValue("nama")
		jk := r.FormValue("jk")
		tempatlahir := r.FormValue("tempatlahir")
		tgllahir := r.FormValue("tgllahir")
		alamat := r.FormValue("alamat")
		tglmasuk := r.FormValue("tglmasuk")
		komplek := r.FormValue("komplek")
		namaayah := r.FormValue("namaayah")
		namaibu := r.FormValue("namaibu")

		jsonData := map[string]string{
			"nis":         nis,
			"nama":        namasantri,
			"jk":          jk,
			"tempatlahir": tempatlahir,
			"tgllahir":    tgllahir,
			"alamat":      alamat,
			"tglmasuk":    tglmasuk,
			"komplek":     komplek,
			"nmayah":      namaayah,
			"nmibu":       namaibu}

		jsonValue, _ := json.Marshal(jsonData)
		client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
		request, err := http.NewRequest("PUT", "http://localhost:8077/api/munawwir/santri/"+newId, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Error :  Tidak dapat memanggil API UPdate Data Santri")
			return
		}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Cannot do request")
			return
		}
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Success Update Data Id :", newId)
		fmt.Println(string(data))
		defer response.Body.Close()
		http.Redirect(w, r, "/santri/dataAll-daftar", 301)
	}
}

func DeleteDataSantri(w http.ResponseWriter, r *http.Request) {
	newId := r.FormValue("id")
	client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
	request, err := http.NewRequest("DELETE", "http://localhost:8077/api/munawwir/santri/"+newId, nil)
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
	if newId == "0" {
		http.Redirect(w, r, "/santri/dataAll-daftar", 301)
	} else {
		http.Redirect(w, r, "/santri/dataAll-daftar", 301)
		fmt.Println("Sukses Hapus Data Santri ID :", newId)
	}
}

func GetAllDataSantriByKomplek(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		var komplek string
		komplekname := r.FormValue("komplek")
		switch komplekname {
		case "110":
			komplek = "Komplek AB"
		case "701":
			komplek = "Komplek Gipa"
		case "081":
			komplek = "Huffadz 1"
		case "082":
			komplek = "Huffadz 2"
		case "901":
			komplek = "Komplek IJ"
		case "012":
			komplek = "Komplek L Putra"
		case "013":
			komplek = "Komplek L Putri"
		case "015":
			komplek = "Komplek M"
		case "231":
			komplek = "Nurussalam Putra"
		case "232":
			komplek = "Nurussalam Putri"
		case "601":
			komplek = "Padang Jagad"
		case "005":
			komplek = "Komplek Q"
		case "061":
			komplek = "Komplek R1"
		case "062":
			komplek = "Komplek R2"
		case "056":
			komplek = "Ribatul Qur'an"
		default:
			komplek = ""
		}
		request, err := http.Get("http://localhost:8077/api/munawwir/santri/komplek/" + komplek)
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var santri []Santri
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &santri)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		var ranges = len(santri)
		fmt.Println("Sukses Ambil Data Santri, Data =", ranges)
		t, err := template.ParseFiles("website/templates/santri/santriall.html")
		if err != nil {
			fmt.Println("Tidak bisa load template Daftar Santri")
		}
		t.Execute(w, santri)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
