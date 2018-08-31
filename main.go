package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ahmadqo/SP_FrontEnd/handler"
	"github.com/ahmadqo/SP_Santri/app/database"
	_ "github.com/go-sql-driver/mysql"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("website/templates/index.html")
	if err != nil {
		fmt.Println("error")
	}
	t.Execute(w, nil)
}
func sejarah(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("website/templates/sejarah.html")
	if err != nil {
		fmt.Println("error")
	}
	t.Execute(w, nil)
}

func DBinit() {
	DB, err := database.NewConnect()
	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connect to Database")
}
func main() {
	DBinit()
	http.HandleFunc("/", index)
	http.HandleFunc("/sejarah", sejarah)

	//DATA SANTRI
	http.HandleFunc("/santri/index", handler.HalamanSantri)
	http.HandleFunc("/santri/dataAll-daftar", handler.GetAllDataSantri)
	http.HandleFunc("/santri/dataInsert-handle", handler.InputDataSantri)
	http.HandleFunc("/santri/data/getSantri-Nis", handler.GetDataSantriByNis)
	http.HandleFunc("/santri/delete/data", handler.DeleteDataSantri)
	http.HandleFunc("/santri/data/update", handler.UpdateDataSantri)
	http.HandleFunc("/santri/data/komplek", handler.GetAllDataSantriByKomplek)

	// NILAI
	http.HandleFunc("/nilai/index", handler.HandleNilaiIndex)
	http.HandleFunc("/nilai/inputData-nilai", handler.HandeInputDataNilai)
	http.HandleFunc("/nilai/tampil/transkipsantri", handler.GetTranskipNilaiSantri)
	http.HandleFunc("/nilai/update/data-by", handler.UpdateNilaiSantri)
	http.HandleFunc("/nilai/deletebykode", handler.DeleteNilaiByKode)

	// JADWAL KAJIAN
	http.HandleFunc("/kajian", handler.HalamanKajian)
	http.HandleFunc("/kajian/data/all", handler.GetDataAllJadwalKajian)
	http.HandleFunc("/kajian/data/id", handler.GetDataByIdJadwalKajian)
	http.HandleFunc("/kajian/data/insert", handler.InsertDataKajian)
	http.HandleFunc("/kajian/data/delete", handler.DeleteDataKajian)
	http.HandleFunc("/kajian/data/update", handler.UpdateDataKajian)

	//JADWAL KEGIATAN
	http.HandleFunc("/kegiatan/santri", handler.HalamanKegiatan)
	http.HandleFunc("/kegiatan/santri/dataJadwal", handler.GetaAllDataKagiatan)
	http.HandleFunc("/kegiatan/santri/dataByid/get", handler.GetByIdDataKegiatan)
	http.HandleFunc("/kegiatan/santri/deleteData", handler.UpdateDataKegiatan)
	http.HandleFunc("/kegiatan/santri/data/insertKegiatan", handler.InsertDataKegiatan)

	http.HandleFunc("/kegiatan/santri/harian", handler.GetaAllDataKagiatanHarian)
	http.HandleFunc("/kegiatan/santri/harian/insert", handler.InsertDataKegiatanHarian)
	http.HandleFunc("/kegiatan/santri/harian/getid", handler.GetByIdDataKegiatanHarian)
	http.HandleFunc("/kegiatan/santri/harian/udpdate", handler.UpdateDataKegiatanHarian)
	http.HandleFunc("/kegiatan/santri/harian/delete", handler.DeleteDataKegiatanHarian)

	// JURNAL
	const uploadPath = "./file"
	http.HandleFunc("/jurnal/index", handler.ShowJurnal)
	http.HandleFunc("/jurnal/action/upload", handler.UploadJurnal)
	http.HandleFunc("/jurnal/delete", handler.DeleteJurnal)

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/jurnal/files/", http.StripPrefix("/jurnal/files", fs))

	// LOGON
	http.HandleFunc("/logon/data", handler.GetDogonData)
	http.HandleFunc("/logon/delete", handler.DeleteLogon)
	http.HandleFunc("/logon/ubahpwadmin", handler.GantiPasswordAdmin)

	// LOGIN USER
	http.HandleFunc("/login", handler.LoginPage)
	http.HandleFunc("/signup", handler.SignUp)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/logout", handler.LogoutUser)

	// USER SANTRI
	http.HandleFunc("/portal/index", handler.OpenIndexSantri)
	http.HandleFunc("/portal/kegiatan/getdata", handler.UserGetDataKegiatan)
	http.HandleFunc("/portal/kajian/getdata", handler.UserGetDataKajian)
	http.HandleFunc("/portal/santri/getdata", handler.UserGetDataSantriInLogin)
	http.HandleFunc("/portal/nilai/getdata", handler.UserGetTranskipNilaiSantri)
	http.HandleFunc("/portal/kegiatan/harian/getdata", handler.UserGetDataKegiatanHarian)
	http.HandleFunc("/portal/jurnal/getdata", handler.UserJurnal)
	http.HandleFunc("/portal/santri/ubahpassword", handler.UserGantiPassword)

	http.Handle("/website/", http.StripPrefix("/website/", http.FileServer(http.Dir("website"))))
	port := ":4020"
	log.Println("Server run on  ", port)
	http.ListenAndServe(port, nil)
}
