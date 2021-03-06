package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"encoding/json"
	"html/template"
	. "github.com/koind/go-photostock/upload"
	"fmt"
	"strconv"
	"time"
)

type WebService struct {
	Uploader
}

func (imgC *WebService) indexAction(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	tmpl.Execute(w, nil)
}

func (imgC *WebService) newAction(w http.ResponseWriter, r *http.Request) {
	uploader := imgC.Uploader

	var imagesName map[int]string
	file, header := uploader.GetFile(r, "image")

	imageType  := uploader.GetImageType(header.Filename)
	imageName  := strconv.FormatInt(time.Now().Unix(), 10) + imageType
	folderPath := "storage/images/"
	imagePath  := folderPath + imageName

	uploader.MkDir(folderPath)
	uploader.MoveFile(file, imagePath)

	imagesName = uploader.DivideByFour(imagePath, folderPath)
	fmt.Println(imagesName)

	if uploader.GetError() != nil {
		http.Error(w, uploader.GetError().Error(), 500)
		return
	}
}

func (imgC *WebService) getOneByIdAction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func (imgC *WebService) getOneByIdAndPartIdAction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func main() {
	ws := WebService{
		Uploader{},
	}
	router := mux.NewRouter()
	router.HandleFunc("/", ws.indexAction).Methods("GET")
	router.HandleFunc("/images", ws.newAction).Methods("POST")
	router.HandleFunc("/images/{id}", ws.getOneByIdAction).Methods("GET")
	router.HandleFunc("/images/{id}/parts/{part_num}", ws.getOneByIdAndPartIdAction).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
