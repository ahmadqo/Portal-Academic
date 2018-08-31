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

type Kajian struct {
	ID                 int    `json:"id"`
	NamaKajian         string `json:"namakajian"`
	HariPelaksanaan    string `json:"haripelaksanaan"`
	TanggalPelaksanaan string `json:"tanggalpelaksanaan"`
	WaktuPelaksanaan   string `json:"waktupelaksanaan"`
	LokasiKajian       string `json:"lokasikajian"`
	Pemateri           string `json:"pemateri"`
}

func HalamanKajian(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		t, err := template.ParseFiles("website/templates/kajian/kajian.html")
		if err != nil {
			fmt.Println("error halaman kajian")
		}
		t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func InsertDataKajian(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		if r.Method != "POST" {
			http.ServeFile(w, r, "website/templates/kajian/kajianinsert.html")
			return
		}
		namakajian := r.FormValue("nama")
		hari := r.FormValue("hari")
		tanggal := r.FormValue("tanggal")
		waktu := r.FormValue("waktu")
		lokasi := r.FormValue("lokasi")
		pemateri := r.FormValue("pemateri")

		jsonData := map[string]string{
			"namakajian":         namakajian,
			"haripelaksanaan":    hari,
			"tanggalpelaksanaan": tanggal,
			"waktupelaksanaan":   waktu,
			"lokasikajian":       lokasi,
			"pemateri":           pemateri}
		jsonValue, _ := json.Marshal(jsonData)
		request, err := http.Post("http://localhost:8088/api/santri/kajian/data", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("The http request failed")
		} else {
			data, _ := ioutil.ReadAll(request.Body)
			fmt.Println(string(data))
		}
		http.Redirect(w, r, "/kajian/data/all", 301)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func UpdateDataKajian(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		var getid = r.FormValue("id")
		request, err := http.Get("http://localhost:8088/api/santri/kajian/" + getid)
		if err != nil {
			fmt.Println(err.Error())
		}
		responseData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		var kajian Kajian
		var jsonData = []byte(responseData)
		err = json.Unmarshal(jsonData, &kajian)
		if err != nil {
			fmt.Println("Cannot Unmarshal data")
		}
		if getid == "0" {
			http.Redirect(w, r, "/kajian/data/all", 301)
		} else {
			t, err := template.ParseFiles("website/templates/kajian/kajianupdate.html")
			if err != nil {
				fmt.Println("Tidak bisa load template kajian update")
			}
			t.Execute(w, kajian)
			return
		}
	} else {
		var getid = r.FormValue("id")
		namakajian := r.FormValue("nama")
		hari := r.FormValue("hari")
		tanggal := r.FormValue("tanggal")
		waktu := r.FormValue("waktu")
		lokasi := r.FormValue("lokasi")
		pemateri := r.FormValue("pemateri")

		jsonData := map[string]string{
			"namakajian":         namakajian,
			"haripelaksanaan":    hari,
			"tanggalpelaksanaan": tanggal,
			"waktupelaksanaan":   waktu,
			"lokasikajian":       lokasi,
			"pemateri":           pemateri}

		jsonValue, _ := json.Marshal(jsonData)
		client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
		request, err := http.NewRequest("PUT", "http://localhost:8088/api/santri/kajian/"+getid, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Cannot call api Updte")
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
		http.Redirect(w, r, "/kajian/data/all", 301)
	}
}
func GetDataAllJadwalKajian(w http.ResponseWriter, r *http.Request) {
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
		t, err := template.ParseFiles("website/templates/kajian/kajianalldata.html")
		if err != nil {
			fmt.Println("Tidak bisa load template kajian all data")
		}
		t.Execute(w, kajians)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func GetDataByIdJadwalKajian(w http.ResponseWriter, r *http.Request) {
	var getid = r.FormValue("id")
	request, err := http.Get("http://localhost:8088/api/santri/kajian/" + getid)
	if err != nil {
		fmt.Println(err.Error())
	}
	responseData, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}
	var kajian Kajian
	var jsonData = []byte(responseData)
	err = json.Unmarshal(jsonData, &kajian)
	if err != nil {
		fmt.Println("Cannot Unmarshal data")
	}
	t, err := template.ParseFiles("website/templates/kajian/kajianbyid.html")
	if err != nil {
		fmt.Println("Tidak bisa load template kajian by id")
	}
	t.Execute(w, kajian)
}
func DeleteDataKajian(w http.ResponseWriter, r *http.Request) {
	newId := r.FormValue("id")
	client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
	request, err := http.NewRequest("DELETE", "http://localhost:8088/api/santri/kajian/"+newId, nil)
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
	http.Redirect(w, r, "/kajian/data/all", 301)
	fmt.Println("Succes Delete :", newId)
}
