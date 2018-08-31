package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const uploadPath = "./file"

type Jurnal struct {
	Kode     string `json:"kode"`
	Namafile string `json:"namafile"`
}

func ShowJurnal(w http.ResponseWriter, r *http.Request) {
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
		t, err := template.ParseFiles("website/templates/jurnal/jurnal.html")
		if err != nil {
			fmt.Println("Tidak Bisa load halaman")
		}
		t.Execute(w, jurs)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func UploadJurnal(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		if r.Method != "POST" {
			http.ServeFile(w, r, "website/templates/jurnal/uploadjurnal.html")
			return
		} else {
			filename := r.FormValue("filename")
			fileType := r.PostFormValue("type")
			file, _, err := r.FormFile("uploadFile")
			if err != nil {
				renderError(w, "INVALID_FILE", http.StatusBadRequest)
				return
			}
			defer file.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Fprintf(w, "<script>alert('Infalid File!');window.location='/jurnal/action/upload'</script>")
				return
			}

			// check file type, detectcontenttype only needs the first 512 bytes
			filetype := http.DetectContentType(fileBytes)
			switch filetype {
			case "application/pdf":
				break
			default:
				fmt.Fprintf(w, "<script>alert('Infalid File Type!');window.location='/jurnal/action/upload'</script>")
				fmt.Println("type file : ", filetype)
				return
			}
			// fileName := randToken(12)
			fileEndings, err := mime.ExtensionsByType(fileType)
			if err != nil {
				fmt.Fprintf(w, "<script>alert('Cannot Read File Type!');window.location='/jurnal/action/upload'</script>")
				return
			}
			newFileName := filename + fileEndings[0]
			newPath := filepath.Join(uploadPath, newFileName)
			fmt.Printf("FileName: %s, FileType: %s, File: %s\n", filename, fileType, newPath)

			// write file
			newFile, err := os.Create(newPath)
			if err != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}
			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}
			jsonData := map[string]string{
				"namafile": newFileName,
			}
			jsonValue, _ := json.Marshal(jsonData)
			request, err := http.Post("http://localhost:8098/jurnal/upload", "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				fmt.Println("The http request filed")
			} else {
				data, _ := ioutil.ReadAll(request.Body)
				fmt.Println(string(data))
				fmt.Println("Succes Insert Jurnal")
				fmt.Fprintf(w, "<script>alert('Berhasil Upload Jurnal!');window.location='/jurnal/action/upload'</script>")
			}
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func DeleteJurnal(w http.ResponseWriter, r *http.Request) {
	newId := r.FormValue("kodejur")
	namafile := r.FormValue("namafile")
	client := &http.Client{} //Statement &http.Client{} menghasilkan instance http.Client
	request, err := http.NewRequest("DELETE", "http://localhost:8098/jurnal/delete/"+newId, nil)
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
	err = os.Remove("file/" + namafile)
	if err != nil {
		fmt.Println("Tidak Bisa menghapus file dalam folder")
	}
	http.Redirect(w, r, "/jurnal/index", 301)
	fmt.Println("Succes Delete :", newId)
}
