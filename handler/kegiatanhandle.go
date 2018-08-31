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

type Kegiatan struct {
	ID                 int    `json:"id"`
	NamaKegiatan       string `json:"namakegiatan"`
	TanggalPelaksanaan string `json:"tanggalpelaksanaan"`
	WaktuPelaksanaan   string `json:"waktupelaksanaan"`
	LokasiKegiatan     string `json:"lokasipelaksanaan"`
	PengampuKegiatan   string `json:"pengampukegiatan"`
	Keterangan         string `json:"keterangan"`
}
type KegiatanHarian struct {
	Id                     int    `json:"idkh"`
	Waktu                  string `json:"waktu"`
	KegiatanMadarasah      string `json:"kegiatanmdrasah"`
	KegiatanTahfidzulQuran string `json:"kegiatantahfidzulwquran"`
	KeteranganKh           string `json:"keterangankh"`
}

func HalamanKegiatan(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		t, err := template.ParseFiles("website/templates/kegiatan/kegiatan.html")
		if err != nil {
			fmt.Println("Error : Tidak Dapat me-Load Halaman Index Kegiatan Santri")
		}
		t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func GetaAllDataKagiatan(w http.ResponseWriter, r *http.Request) {
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
		t, err := template.ParseFiles("website/templates/kegiatan/kegiatanalldata.html")
		if err != nil {
			fmt.Println("Tidak bisa load template Daftar Kegiatan")
		}
		t.Execute(w, kegiatan)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func GetByIdDataKegiatan(w http.ResponseWriter, r *http.Request) {
	var getid = r.FormValue("id")
	request, err := http.Get("http://localhost:8090/api/santri/kegiatan/" + getid)
	if err != nil {
		fmt.Println(err.Error())
	}
	responseData, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}
	var kegiatan Kegiatan
	var jsonData = []byte(responseData)
	err = json.Unmarshal(jsonData, &kegiatan)
	if err != nil {
		fmt.Println("Cannot Unmarshal data")
	}
	t, err := template.ParseFiles("website/templates/kegiatan/kegiatanbyid.html")
	if err != nil {
		fmt.Println("Tidak bisa load template Detil Kegiatan")
	}
	t.Execute(w, kegiatan)
}

func InsertDataKegiatan(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		if r.Method != "POST" {
			http.ServeFile(w, r, "website/templates/kegiatan/kegiataninsert.html")
			return
		}
		namakajian := r.FormValue("nama")
		tanggal := r.FormValue("tanggal")
		waktu := r.FormValue("waktu")
		lokasi := r.FormValue("lokasi")
		pengampu := r.FormValue("pengampu")
		keterangan := r.FormValue("keterangan")

		jsonData := map[string]string{
			"namakegiatan":       namakajian,
			"tanggalpelaksanaan": tanggal,
			"waktupelaksanaan":   waktu,
			"lokasipelaksanaan":  lokasi,
			"pengampukegiatan":   pengampu,
			"keterangan":         keterangan}

		jsonValue, _ := json.Marshal(jsonData)
		request, err := http.Post("http://localhost:8090/api/santri/kegiatan/data", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("The http request failed")
		} else {
			data, _ := ioutil.ReadAll(request.Body)
			fmt.Println(string(data))
		}
		http.Redirect(w, r, "/kegiatan/santri/dataJadwal", 301)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func UpdateDataKegiatan(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		var getid = r.FormValue("id")
		request, err := http.Get("http://localhost:8090/api/santri/kegiatan/" + getid)
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var kegiatan Kegiatan
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &kegiatan)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		if getid == "0" {
			http.Redirect(w, r, "/kegiatan/santri/dataJadwal", 301)
		} else {
			t, err := template.ParseFiles("website/templates/kegiatan/kegiatanupdate.html")
			if err != nil {
				fmt.Println("Tidak Dapat me-Load Halaman Update Kegiatan")
			}
			t.Execute(w, kegiatan)
			return
		}
	} else {
		var getid = r.FormValue("id")
		namakajian := r.FormValue("nama")
		tanggal := r.FormValue("tanggal")
		waktu := r.FormValue("waktu")
		lokasi := r.FormValue("lokasi")
		pengampu := r.FormValue("pengampu")
		keterangan := r.FormValue("keterangan")

		jsonData := map[string]string{
			"namakegiatan":       namakajian,
			"tanggalpelaksanaan": tanggal,
			"waktupelaksanaan":   waktu,
			"lokasipelaksanaan":  lokasi,
			"pengampukegiatan":   pengampu,
			"keterangan":         keterangan}

		jsonValue, _ := json.Marshal(jsonData)
		client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
		request, err := http.NewRequest("PUT", "http://localhost:8090/api/santri/kegiatan/"+getid, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Error :  Tidak dapat memanggil API UPdate Data Kegiatan")
			return
		}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Cannot do request")
			return
		}
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Success Update Data Id :", getid)
		fmt.Println(string(data))
		defer response.Body.Close()
		http.Redirect(w, r, "/kegiatan/santri/dataJadwal", 301)
	}
}

func DeleteDataKegiatanById(w http.ResponseWriter, r *http.Request) {
	newId := r.FormValue("id")
	client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
	request, err := http.NewRequest("DELETE", "http://localhost:8090/api/santri/kegiatan/"+newId, nil)
	if err != nil {
		fmt.Println("Error : Tidak bisa memanggil API Delete")
		return
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Cannot do request")
		return
	}
	defer response.Body.Close()
	http.Redirect(w, r, "/kegiatan/santri/dataJadwal", 301)
	fmt.Println("Succes Delete :", newId)
}

// ============================================================================================
func GetaAllDataKagiatanHarian(w http.ResponseWriter, r *http.Request) {
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
		t, err := template.ParseFiles("website/templates/kegiatan/kegiatanharianall.html")
		if err != nil {
			fmt.Println("Tidak bisa load template Daftar Kegiatan")
		}
		t.Execute(w, kegharian)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func InsertDataKegiatanHarian(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		if r.Method != "POST" {
			http.ServeFile(w, r, "website/templates/kegiatan/kegiatanharianinsert.html")
			return
		}
		waktu := r.FormValue("waktu")
		madrasah := r.FormValue("madrasah")
		tahfidzul := r.FormValue("tahfidzul")
		keterangan := r.FormValue("keterangan")

		jsonData := map[string]string{
			"waktu":                   waktu,
			"kegiatanmdrasah":         madrasah,
			"kegiatantahfidzulwquran": tahfidzul,
			"keterangankh":            keterangan}

		jsonValue, _ := json.Marshal(jsonData)
		request, err := http.Post("http://localhost:8090/api/kegiatan/harian/santri/insert", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("The http request failed")
		} else {
			data, _ := ioutil.ReadAll(request.Body)
			fmt.Println(string(data))
		}
		http.Redirect(w, r, "/kegiatan/santri/harian", 301)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func GetByIdDataKegiatanHarian(w http.ResponseWriter, r *http.Request) {
	var getid = r.FormValue("id")
	request, err := http.Get("http://localhost:8090/api/kegiatan/harian/santri/" + getid)
	if err != nil {
		fmt.Println(err.Error())
	}
	responseData, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}
	var kegharian KegiatanHarian
	var jsonData = []byte(responseData)
	err = json.Unmarshal(jsonData, &kegharian)
	if err != nil {
		fmt.Println("Cannot Unmarshal data")
	}
	t, err := template.ParseFiles("website/templates/kegiatan/kegiatanharianbyid.html")
	if err != nil {
		fmt.Println("Tidak bisa load template Detil Kegiatan")
	}
	t.Execute(w, kegharian)
}
func UpdateDataKegiatanHarian(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		var getid = r.FormValue("id")
		request, err := http.Get("http://localhost:8090/api/kegiatan/harian/santri/" + getid)
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var kegharian KegiatanHarian
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &kegharian)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		if getid == "0" {
			http.Redirect(w, r, "/kegiatan/santri/dataJadwal", 301)
		} else {
			t, err := template.ParseFiles("website/templates/kegiatan/kegiatanharianupdate.html")
			if err != nil {
				fmt.Println("Tidak Dapat me-Load Halaman Update Kegiatan")
			}
			t.Execute(w, kegharian)
			return
		}
	} else {
		var getid = r.FormValue("id")
		waktu := r.FormValue("waktu")
		madrasah := r.FormValue("madrasah")
		tahfidzul := r.FormValue("tahfidzul")
		keterangan := r.FormValue("keterangan")

		jsonData := map[string]string{
			"waktu":                   waktu,
			"kegiatanmdrasah":         madrasah,
			"kegiatantahfidzulwquran": tahfidzul,
			"keterangankh":            keterangan}

		jsonValue, _ := json.Marshal(jsonData)
		client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
		request, err := http.NewRequest("PUT", "http://localhost:8090/api/kegiatan/harian/santri/"+getid, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Error :  Tidak dapat memanggil API UPdate Data Kegiatan")
			return
		}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Cannot do request")
			return
		}
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Success Update Data Id :", getid)
		fmt.Println(string(data))
		defer response.Body.Close()
		http.Redirect(w, r, "/kegiatan/santri/harian", 301)
	}
}

func DeleteDataKegiatanHarian(w http.ResponseWriter, r *http.Request) {
	newId := r.FormValue("id")
	client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
	request, err := http.NewRequest("DELETE", "http://localhost:8090/api/kegiatan/harian/santri/"+newId, nil)
	if err != nil {
		fmt.Println("Error : Tidak bisa memanggil API Delete")
		return
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Cannot do request")
		return
	}
	defer response.Body.Close()
	http.Redirect(w, r, "/kegiatan/santri/harian", 301)
	fmt.Println("Succes Delete :", newId)
}
